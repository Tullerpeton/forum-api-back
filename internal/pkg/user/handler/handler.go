package handler

import (
	"encoding/json"
	"net/http"

	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/user"
	"github.com/forum-api-back/pkg/tools/http_utils"
	"github.com/forum-api-back/pkg/errors"

	"github.com/valyala/fasthttp"
	"github.com/fasthttp/router"
)

type UserHandler struct {
	UserUCase user.UseCase
}

func NewHandler(userUCase user.UseCase) user.Handler {
	return &UserHandler{
		UserUCase: userUCase,
	}
}

func (h *UserHandler) InitHandler(r *router.Router) {
	r.POST("/api/user/{nickname}/create", h.CreateNewUser)
	r.GET("/api/user/{nickname}/profile", h.GetUserProfile)
	r.POST("/api/user/{nickname}/profile", h.UpdateUserProfile)
}

func (h *UserHandler) CreateNewUser(ctx *fasthttp.RequestCtx) {
	userInfo := &models.User{}
	if err := json.Unmarshal(ctx.PostBody(), userInfo); err != nil {
		http_utils.SetJSONResponse(ctx, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if userInfo.NickName = ctx.UserValue("nickname").(string); userInfo.NickName == "" {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	err := h.UserUCase

}

func (h *UserHandler) GetUserProfile(ctx *fasthttp.RequestCtx) {

}

func (h *UserHandler) UpdateUserProfile(ctx *fasthttp.RequestCtx) {

}