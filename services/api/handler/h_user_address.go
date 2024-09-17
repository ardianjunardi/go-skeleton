package handler

import (
	"context"
	"go-skeleton/bootstrap"
	"go-skeleton/lib/utils"
	"go-skeleton/services/api/model"
	"go-skeleton/services/api/request"
	"go-skeleton/services/api/response"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Contract) GetAddressByAddressIdentifier(w http.ResponseWriter, r *http.Request) {
	var (
		err               error
		ctx               = context.TODO()
		m                 = model.Contract{App: h.App}
		dataUserAddress   = model.UserAddressEnt{}
		res               = response.UserAddressRes{}
		addressIdentifier = chi.URLParam(r, "code")
	)

	// Fetch user address by address identifier
	dataUserAddress, err = m.GetAddressByAddressIdentifier(h.DB, ctx, addressIdentifier)
	if err != nil {
		// if empty data still success response
		if err.Error() == utils.EmptyData {
			h.SendEmptyDataSuccess(w, res, nil)
			return
		}
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	res = response.UserAddressRes{
		AddressIdentifier: dataUserAddress.AddressIdentifier,
		Title:             dataUserAddress.Title.String,
		FullAddress:       dataUserAddress.FullAddress,
		CreatedDate:       dataUserAddress.CreatedDate.Format(utils.DATE_TIME_FORMAT),
		UpdatedDate:       dataUserAddress.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
		DeletedDate:       dataUserAddress.DeletedDate.Time.Format(utils.DATE_TIME_FORMAT),
	}

	h.SendSuccess(w, res, nil)
}

func (h *Contract) GetAllAddressesByUserIdentifier(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		m              = model.Contract{App: h.App}
		dataUser       = model.UserEnt{}
		res            = []response.UserAddressRes{}
		userIdentifier = bootstrap.GetUserIdentifierFromToken(ctx, r)
	)

	dataUser, err = m.GetUserByUserIdentifier(h.DB, ctx, userIdentifier)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Fetch user address by address identifier
	dataUserAddresses, err := m.GetUserAddressesByUserID(h.DB, ctx, int64(dataUser.ID))
	if err != nil {
		// if empty data still success response
		if err.Error() == utils.EmptyData {
			h.SendEmptyDataSuccess(w, res, nil)
			return
		}
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	for _, dataUserAddress := range dataUserAddresses {
		res = append(res, response.UserAddressRes{
			AddressIdentifier: dataUserAddress.AddressIdentifier,
			Title:             dataUserAddress.Title.String,
			FullAddress:       dataUserAddress.FullAddress,
			CreatedDate:       dataUserAddress.CreatedDate.Format(utils.DATE_TIME_FORMAT),
			UpdatedDate:       dataUserAddress.UpdatedDate.Time.Format(utils.DATE_TIME_FORMAT),
			DeletedDate:       dataUserAddress.DeletedDate.Time.Format(utils.DATE_TIME_FORMAT),
		})
	}

	h.SendSuccess(w, res, nil)
}

func (h *Contract) InsertUserAddressAct(w http.ResponseWriter, r *http.Request) {
	var (
		err            error
		ctx            = context.TODO()
		req            = request.UserAddressReq{}
		m              = model.Contract{App: h.App}
		userIdentifier = bootstrap.GetUserIdentifierFromToken(ctx, r)
		dataUser       = model.UserEnt{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err.Error())
		return
	}

	dataUser, err = m.GetUserByUserIdentifier(h.DB, ctx, userIdentifier)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Generate User Identifier
	addressIdentifier := utils.GeneratePrefixCode(utils.UserAddressPrefix)
	err = m.InsertUserAddress(h.DB, ctx, int64(dataUser.ID), addressIdentifier, req.Title, req.FullAddress)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

func (h *Contract) UpdateUserAddressAct(w http.ResponseWriter, r *http.Request) {
	var (
		err               error
		ctx               = context.TODO()
		req               = request.UserAddressReq{}
		m                 = model.Contract{App: h.App}
		addressIdentifier = chi.URLParam(r, "code")
		userIdentifier    = bootstrap.GetUserIdentifierFromToken(ctx, r)
		dataUser          = model.UserEnt{}
	)

	// Bind and validate
	if err = h.BindAndValidate(r, &req); err != nil {
		h.SendBindAndValidateError(w, err.Error())
		return
	}

	dataUser, err = m.GetUserByUserIdentifier(h.DB, ctx, userIdentifier)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.UpdateUserAddress(h.DB, ctx, int64(dataUser.ID), addressIdentifier, req.Title, req.FullAddress)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}

func (h *Contract) DeleteUserAddressAct(w http.ResponseWriter, r *http.Request) {
	var (
		err               error
		ctx               = context.TODO()
		m                 = model.Contract{App: h.App}
		addressIdentifier = chi.URLParam(r, "code")
		userIdentifier    = bootstrap.GetUserIdentifierFromToken(ctx, r)
		dataUser          = model.UserEnt{}
	)

	dataUser, err = m.GetUserByUserIdentifier(h.DB, ctx, userIdentifier)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	err = m.DeleteUserAddress(h.DB, ctx, int64(dataUser.ID), addressIdentifier)
	if err != nil {
		h.SendBadRequest(w, err.Error())
		return
	}

	// Populate response
	h.SendSuccess(w, nil, nil)
}
