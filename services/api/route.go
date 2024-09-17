package api

import (
	"go-skeleton/bootstrap"
	"go-skeleton/services/api/handler"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes all routes for the apps
func RegisterRoutes(r *chi.Mux, app *bootstrap.App) {
	r.Route("/v1", func(r chi.Router) {
		r.Get("/ping", app.PingAction)

		AppSubsRoute(r, app)
		// CMSSubsRoute(r, app)
	})
}

func AppSubsRoute(r chi.Router, app *bootstrap.App) {
	h := handler.Contract{App: app}

	r.Route("/auths", func(r chi.Router) {
		r.Post("/login", h.LoginUserAct)
		r.Post("/register", h.RegisterUserAct)

		// Request Token for Registration and ForgotPassword
		r.Post("/request-token", h.RequestVerifyEmailUserAct)
		r.Post("/verify-token", h.VerifyTokenUserAct)
		r.With(app.VerifyJwtTokenUser).Post("/reset-password", h.ResetPasswordUserAct)
	})

	// User
	r.Route("/users", func(r chi.Router) {
		// User Profile
		r.Route("/profile", func(r chi.Router) {
			r.Use(app.VerifyJwtTokenUser)
			r.Get("/", h.GetUserProfileAct)
			r.Put("/", h.UpdateUserProfileAct)
		})

		// User Addresses
		r.Route("/addresses", func(r chi.Router) {
			r.Use(app.VerifyJwtTokenUser)
			r.Get("/", h.GetAllAddressesByUserIdentifier)
			r.Get("/{code}", h.GetAddressByAddressIdentifier)
			r.Post("/", h.InsertUserAddressAct)
			r.Put("/{code}", h.UpdateUserAddressAct)
			r.Delete("/{code}", h.DeleteUserAddressAct)
		})

		r.With(app.VerifyJwtTokenUser).Put("/update-password", h.UpdatePasswordUserAct)
	})

	// Master Setting
	r.Route("/settings", func(r chi.Router) {
		r.Get("/", h.GetSettingListAct)
		r.Get("/{code}", h.GetSettingDetailAct)
		r.With(app.VerifyJwtTokenUser).Post("/", h.AddSettingAct)
		r.With(app.VerifyJwtTokenUser).Put("/{code}", h.UpdateSettingAct)
	})

	// Upload
	r.Route("/uploads", func(r chi.Router) {
		r.Use(app.VerifyJwtTokenUser)
		r.Post("/", h.UploadFileAct)
	})
}
