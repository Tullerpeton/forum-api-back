package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/user"
	"github.com/forum-api-back/pkg/errors"
	"github.com/forum-api-back/pkg/tools/http_utils"

	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	UserUCase user.UseCase
}

func NewHandler(userUCase user.UseCase) user.Handler {
	return &UserHandler{
		UserUCase: userUCase,
	}
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

	newUser, err := h.UserUCase.CreateNewUser(userInfo)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, newUser[0], http.StatusCreated)
	case errors.ErrAlreadyExists:
		http_utils.SetJSONResponse(ctx, newUser, http.StatusConflict)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetUserProfile(ctx *fasthttp.RequestCtx) {
	userNickName := ctx.UserValue("nickname").(string)
	if userNickName == "" {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	selectedUser, err := h.UserUCase.GetUserByNickName(userNickName)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, selectedUser, http.StatusOK)
	case errors.ErrUserNotFound:
		http_utils.SetJSONResponse(ctx, errors.ErrUserNotFound, http.StatusNotFound)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetUsersByForum(ctx *fasthttp.RequestCtx) {
	forumSlug := ctx.UserValue("slug").(string)
	if forumSlug == "" {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	userPaginator := &models.UserPaginator{Limit: 100}
	userPaginator.Since = string(ctx.FormValue("since"))

	if isDesc := string(ctx.FormValue("desc")); isDesc == "true" {
		userPaginator.SortOrder = true
	}

	parseLimit, err := strconv.Atoi(string(ctx.FormValue("limit")))
	if err == nil {
		userPaginator.Limit = uint64(parseLimit)
	}

	selectedUsers, err := h.UserUCase.GetUsersByForum(forumSlug, userPaginator)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, selectedUsers, http.StatusOK)
	case errors.ErrForumNotFound:
		http_utils.SetJSONResponse(ctx, errors.Error{Message: fmt.Sprintf("Can't find forum by slug: %s", forumSlug)}, http.StatusNotFound)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}

func (h *UserHandler) UpdateUserProfile(ctx *fasthttp.RequestCtx) {
	userInfo := &models.User{}
	if err := json.Unmarshal(ctx.PostBody(), userInfo); err != nil {
		http_utils.SetJSONResponse(ctx, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	if userInfo.NickName = ctx.UserValue("nickname").(string); userInfo.NickName == "" {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	updatedUser, err := h.UserUCase.SetUserProfile(userInfo)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, updatedUser, http.StatusOK)
	case errors.ErrUserNotFound:
		http_utils.SetJSONResponse(ctx, errors.ErrUserNotFound, http.StatusNotFound)
	case errors.ErrAlreadyExists:
		http_utils.SetJSONResponse(ctx, errors.ErrAlreadyExists, http.StatusConflict)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}
