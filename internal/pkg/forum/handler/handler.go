package handler

import (
	"encoding/json"
	"net/http"

	"github.com/forum-api-back/internal/pkg/forum"
	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/pkg/errors"
	"github.com/forum-api-back/pkg/tools/http_utils"

	"github.com/valyala/fasthttp"
)

type ForumHandler struct {
	ForumUCase forum.UseCase
}

func NewHandler(forumUCase forum.UseCase) forum.Handler {
	return &ForumHandler{
		ForumUCase: forumUCase,
	}
}

func (h *ForumHandler) CreateNewForum(ctx *fasthttp.RequestCtx) {
	forumInfo := &models.ForumCreate{}
	if err := json.Unmarshal(ctx.PostBody(), forumInfo); err != nil {
		http_utils.SetJSONResponse(ctx, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	newForum, err := h.ForumUCase.CreateNewForum(forumInfo)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, newForum, http.StatusCreated)
	case errors.ErrUserNotFound:
		http_utils.SetJSONResponse(ctx, errors.ErrUserNotFound, http.StatusNotFound)
	case errors.ErrAlreadyExists:
		http_utils.SetJSONResponse(ctx, newForum, http.StatusConflict)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}

func (h *ForumHandler) GetForumDetails(ctx *fasthttp.RequestCtx) {
	forumSlug := ctx.UserValue("slug").(string)
	if forumSlug == "" {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	selectedForum, err := h.ForumUCase.GetForumDetails(forumSlug)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, selectedForum, http.StatusOK)
	case errors.ErrForumNotFound:
		http_utils.SetJSONResponse(ctx, errors.ErrForumNotFound, http.StatusNotFound)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}
