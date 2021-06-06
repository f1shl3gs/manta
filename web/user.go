package web

import (
	"encoding/json"
	"net/http"

	"github.com/f1shl3gs/manta"
	"github.com/f1shl3gs/manta/authz"
	"go.uber.org/zap"
)

const (
	userPrefix = "/api/v1/users"
	userIDPath = "/api/v1/users/:id"
	viewerPath = "/api/v1/viewer"
)

type UserHandler struct {
	*Router

	logger      *zap.Logger
	userService manta.UserService
}

func userService(logger *zap.Logger, router *Router, svc manta.UserService) {
	h := &UserHandler{
		Router:      router,
		logger:      logger,
		userService: svc,
	}

	h.HandlerFunc(http.MethodGet, viewerPath, h.viewerHandler)
	h.HandlerFunc(http.MethodPost, userPrefix, h.handleAdd)
	h.HandlerFunc(http.MethodGet, userPrefix, h.handleList)
	h.HandlerFunc(http.MethodDelete, userIDPath, h.handleDelete)
}

func (h *UserHandler) viewerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authorizer := authz.FromContext(ctx)
	u, err := h.userService.FindUserByID(ctx, authorizer.GetUserID())
	if err != nil {
		h.HandleHTTPError(ctx, err, w)
		return
	}

	err = encodeResponse(r.Context(), w, http.StatusOK, u)
	if err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func (h *UserHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := idFromRequest(r)
	if err != nil {
		err = encodeResponse(ctx, w, http.StatusBadRequest, manta.ErrInvalidID)
		if err != nil {
			logEncodingError(h.logger, r, err)
		}

		return
	}

	err = h.userService.DeleteUser(ctx, id)
	if err != nil {
		err = encodeResponse(ctx, w, http.StatusInternalServerError, manta.Error{
			Code: manta.EInternal,
			Op:   "delete user",
			Err:  err,
		})
		if err != nil {
			logEncodingError(h.logger, r, err)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) handleList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.userService.FindUsers(ctx, manta.UserFilter{})
	if err != nil {
		if err := encodeResponse(ctx, w, http.StatusInternalServerError, err); err != nil {
			logEncodingError(h.logger, r, err)
		}

		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, &users); err != nil {
		logEncodingError(h.logger, r, err)
	}
}

func decodeUser(r *http.Request) (*manta.User, error) {
	u := &manta.User{}
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (h *UserHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := decodeUser(r)
	if err != nil {
		err = encodeResponse(ctx, w, http.StatusBadRequest, manta.Error{
			Code: manta.EInvalid,
			Op:   "decode",
			Err:  err,
		})
		if err != nil {
			logEncodingError(h.logger, r, err)
		}

		return
	}

	err = h.userService.CreateUser(ctx, user)
	if err != nil {
		err = encodeResponse(ctx, w, http.StatusInternalServerError, manta.Error{
			Code: manta.EInternal,
			Err:  err,
		})
		if err != nil {
			logEncodingError(h.logger, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

/*
func (h *UserHandler) viewerHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	auth := authorizer.FromContext(r.Context())
	if auth == nil {
		encodeResponse(ctx, w, http.StatusUnauthorized, nil)
		return
	}

	if err := encodeResponse(ctx, w, http.StatusOK, auth); err != nil {
		logEncodingError(h.logger, r, err)
	}
}*/
