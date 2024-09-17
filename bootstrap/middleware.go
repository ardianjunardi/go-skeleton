package bootstrap

import (
	"context"
	"fmt"
	"go-skeleton/lib/utils"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type CustomUserClaims struct {
	UserIdentifier string `json:"user_identifier"`
	Email          string `json:"email"`
	jwt.StandardClaims
}

const (
	ChannelApp = "app"
	ChannelCMS = "cms"
)

var (
	mustHeader = []string{"X-Channel", "Content-Type"}
	headerVal  = []string{ChannelApp, ChannelCMS, "application/json"}
)

func userContext(ctx context.Context, subject, id interface{}) context.Context {
	return context.WithValue(ctx, subject, id)
}

const pingReqURI string = "/v1/ping"

func isPingRequest(r *http.Request) bool {
	return r.RequestURI == pingReqURI
}

// NotfoundMiddleware A custom not found response.
func (app *App) NotfoundMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tctx := chi.NewRouteContext()
		rctx := chi.RouteContext(r.Context())

		if !rctx.Routes.Match(tctx, r.Method, r.URL.Path) {
			app.SendNotfound(w, utils.ErrNotFoundPage)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					debug.PrintStack()
				}

				app.Log.FromDefault().WithFields(logrus.Fields{
					"Panic": rvr,
				}).Errorf("Panic: %v \n %v", rvr, string(debug.Stack()))

				app.SendBadRequest(w, utils.ErrSystemError)
				return
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (app *App) VerifyJwtTokenUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := &CustomUserClaims{}
		tokenAuth := r.Header.Get("Authorization")
		_, err := jwt.ParseWithClaims(tokenAuth, claims, func(token *jwt.Token) (interface{}, error) {
			if jwt.SigningMethodHS256 != token.Method {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			secret := app.Config.GetString("app.key")
			return []byte(secret), nil
		})

		if err != nil {
			msg := utils.ErrInvalidToken
			if mErr, ok := err.(*jwt.ValidationError); ok {
				if mErr.Errors == jwt.ValidationErrorExpired {
					msg = utils.ErrTokenExpired
				}
			}

			app.SendAuthError(w, msg)
			return
		}

		// check if token expired or not
		if claims.ExpiresAt <= time.Now().Unix() {
			app.SendAuthError(w, utils.ErrTokenExpired)
			return
		}

		ctx := userContext(r.Context(), "identifier", map[string]string{
			"user_identifier": claims.UserIdentifier,
			"email":           claims.Email,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HeaderCheckerMiddleware check the necesarry headers
func (app *App) HeaderCheckerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range mustHeader {
			if len(r.Header.Get(v)) == 0 || !utils.Contains(headerVal, r.Header.Get(v)) {
				app.SendBadRequest(w, fmt.Sprintf("undefined %s header or wrong value of header", v))
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
