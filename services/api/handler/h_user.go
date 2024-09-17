package handler

import (
	"context"
	"go-skeleton/bootstrap"
	"go-skeleton/lib/utils"
	"go-skeleton/services/api/model"
	"go-skeleton/services/api/request"
	"go-skeleton/services/api/response"
	"net/http"
)

func (h *Contract) GetUserProfileAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		userIdentifier = bootstrap.GetUserIdentifierFromToken(ctx, r)
		m              = model.Contract{App: h.App}
		dataUser       model.UserEnt
		res            = response.UserProfileRes{}
	)

	dataUser, err = m.GetUserByUserIdentifier(h.DB, ctx, userIdentifier)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	res = response.UserProfileRes{
		UserIdentifier: dataUser.UserIdentifier,
		FirstName:      dataUser.FirstName,
		LastName:       dataUser.LastName.String,
		Email:          dataUser.Email,
		AvatarURL:      dataUser.AvatarURL.String,
		IsVerified:     dataUser.IsVerified,
		CreatedDate:    dataUser.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		UpdatedDate:    dataUser.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
	}

	h.SendSuccess(w, res, nil)
}

func (h *Contract) UpdateUserProfileAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		req            = request.UpdateProfileUserReq{}
		m              = model.Contract{App: h.App}
		userIdentifier = bootstrap.GetUserIdentifierFromToken(ctx, r)
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err.Error())
		return
	}

	err = m.UpdateUserProfile(h.DB, ctx, userIdentifier, req.FirstName, req.LastName, req.Description, req.AvatarUrl)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	h.SendSuccess(w, nil, nil)
}

func (h *Contract) UpdatePasswordUserAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		userIdentifier = bootstrap.GetUserIdentifierFromToken(ctx, r)
		m              = model.Contract{App: h.App}
		req            = request.UpdatePasswordReq{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err)
		return
	}

	err = m.UpdatePasswordUser(h.DB, ctx, userIdentifier, req.OldPassword, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}
