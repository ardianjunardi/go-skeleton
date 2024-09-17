package handler

import (
	"context"
	"go-skeleton/bootstrap"
	"go-skeleton/lib/utils"
	"go-skeleton/services/api/model"
	"go-skeleton/services/api/request"
	"go-skeleton/services/api/response"
	"net/http"
	"time"
)

func (h *Contract) LoginUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		ctx = context.TODO()
		m   = model.Contract{App: h.App}
		req = request.LoginUserReq{}
		res = response.LoginUserRes{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	dataUser, expAtUnix, jwtToken, err := m.UserLogin(h.DB, ctx, req.Email, req.Password)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// convert unix timestamp
	expAt := time.Unix(expAtUnix, 0)
	res = response.LoginUserRes{
		Token:          jwtToken,
		AvatarURL:      dataUser.AvatarURL.String,
		UserIdentifier: dataUser.UserIdentifier,
		FirstName:      dataUser.FirstName,
		LastName:       dataUser.LastName.String,
		Email:          dataUser.Email,
		ExpiredAt:      expAt.Format(utils.DATE_TIME_FORMAT),
		ActorType:      utils.User,
		CreatedDate:    time.Now().In(time.UTC).Format(utils.DATE_TIME_FORMAT),
	}
	// Populate Response
	h.SendSuccess(w, res, nil)
}

func (h *Contract) RegisterUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		ctx = context.TODO()
		m   = model.Contract{App: h.App}
		req = request.RegisterUserReq{}
		res = response.RegisterUserRes{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	// Compare Password
	if req.Password != req.ConfirmPassword {
		h.SendBadRequest(w, utils.ErrPasswordMismatch)
		return
	}

	userIdentifier, email, err := m.RegisterUser(h.DB, ctx, req.FirstName, req.LastName, req.Email, req.ConfirmPassword)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.RequestVerifyEmailUser(h.DB, ctx, req.Email, utils.VerifyRegistration)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	res = response.RegisterUserRes{
		UserIdentifier: userIdentifier,
		Email:          email,
		CreatedDate:    time.Now().In(time.UTC).Format(utils.DATE_TIME_FORMAT),
	}

	h.SendSuccess(w, res, nil)
}

func (h *Contract) RequestVerifyEmailUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		check bool
		ctx   = context.TODO()
		m     = model.Contract{App: h.App}
		req   = request.RequestVerifyEmailReq{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	if check = utils.Contains(utils.VerificationType, req.Type); !check {
		h.SendBadRequest(w, utils.ErrInvalidSendingEmailType)
		return
	}

	err = m.RequestVerifyEmailUser(h.DB, ctx, req.Email, req.Type)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

func (h *Contract) VerifyTokenUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		ctx       = context.TODO()
		m         = model.Contract{App: h.App}
		jwtToken  string
		expAt     time.Time
		expAtUnix int64
		dataUser  model.UserEnt
		res       = response.LoginUserRes{}

		// Initiate Query Param
		param = map[string]interface{}{
			"type":  "",
			"token": "",
		}
	)

	if token, ok := r.URL.Query()["token"]; ok && len(token[0]) > 0 {
		param["token"] = token[0]
	}

	if types, ok := r.URL.Query()["type"]; ok && len(types) > 0 {
		paramType := types[0]
		if !utils.Contains(utils.VerificationType, paramType) {
			h.SendBadRequest(w, utils.ErrInvalidSendingEmailType)
			return
		}
		param["type"] = paramType
	} else {
		h.SendBadRequest(w, utils.ErrInvalidTypeQueryParameter)
		return
	}

	dataUser, expAtUnix, jwtToken, err = m.CheckTokenAndExpiration(h.DB, ctx, param["type"].(string), utils.User, param["token"].(string))
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	//convert unix timestamp
	expAt = time.Unix(expAtUnix, 0)

	// Populate response
	res = response.LoginUserRes{
		Token:          jwtToken,
		AvatarURL:      dataUser.AvatarURL.String,
		UserIdentifier: dataUser.UserIdentifier,
		FirstName:      dataUser.FirstName,
		LastName:       dataUser.LastName.String,
		Email:          dataUser.Email,
		ExpiredAt:      expAt.Format(utils.DATE_TIME_FORMAT),
		ActorType:      utils.User,
		CreatedDate:    time.Now().In(time.UTC).Format(utils.DATE_TIME_FORMAT),
	}

	h.SendSuccess(w, res, nil)
}

func (h *Contract) ResetPasswordUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		userIdentifier = bootstrap.GetUserIdentifierFromToken(ctx, r)
		m              = model.Contract{App: h.App}
		req            = request.ResetPasswordReq{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	err = m.ResetPassword(h.DB, ctx, userIdentifier, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}
