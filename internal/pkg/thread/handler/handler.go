package handler

import (
	"encoding/json"
	"net/http"

	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/thread"
	"github.com/forum-api-back/pkg/errors"
	"github.com/forum-api-back/pkg/tools/http_utils"

	"github.com/valyala/fasthttp"
)

type ThreadHandler struct {
	ThreadUCase thread.UseCase
}

func NewHandler(threadUCase thread.UseCase) thread.Handler {
	return &ThreadHandler{
		ThreadUCase: threadUCase,
	}
}

func (h *ThreadHandler) CreateNewThread(ctx *fasthttp.RequestCtx) {
	threadInfo := &models.ThreadCreate{}
	if err := json.Unmarshal(ctx.PostBody(), threadInfo); err != nil {
		http_utils.SetJSONResponse(ctx, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	forumSlug := ctx.UserValue("slug").(string)
	if forumSlug == "" {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	newThread, err := h.ThreadUCase.CreateNewThread(forumSlug, threadInfo)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, newThread, http.StatusCreated)
	case errors.ErrDataConflict:
		http_utils.SetJSONResponse(ctx, errors.ErrDataConflict, http.StatusNotFound)
	case errors.ErrAlreadyExists:
		http_utils.SetJSONResponse(ctx, newThread, http.StatusConflict)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}

func (h *ThreadHandler) GetThreadsByForum(ctx *fasthttp.RequestCtx) {
	threadPaginator := &models.ThreadPaginator{Limit: 100}
	if err := json.Unmarshal(ctx.PostBody(), threadPaginator); err != nil {
		http_utils.SetJSONResponse(ctx, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	forumSlug := ctx.UserValue("slug").(string)
	if forumSlug == "" {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	selectedThreads, err := h.ThreadUCase.GetThreadsByForum(forumSlug, threadPaginator)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, selectedThreads, http.StatusOK)
	case errors.ErrForumNotFound:
		http_utils.SetJSONResponse(ctx, errors.ErrForumNotFound, http.StatusNotFound)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}

func (h *ThreadHandler) GetThreadDetails(ctx *fasthttp.RequestCtx) {
	threadSlugOrId := ctx.UserValue("slug_or_id").(string)
	if threadSlugOrId == "" {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	threadDetails, err := h.ThreadUCase.GetThreadDetails(threadSlugOrId)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, threadDetails, http.StatusOK)
	case errors.ErrThreadNotFound:
		http_utils.SetJSONResponse(ctx, errors.ErrThreadNotFound, http.StatusNotFound)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}

func (h *ThreadHandler) UpdateThreadDetails(ctx *fasthttp.RequestCtx) {
	threadUpdate := &models.ThreadUpdate{}
	if err := json.Unmarshal(ctx.PostBody(), threadUpdate); err != nil {
		http_utils.SetJSONResponse(ctx, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	threadSlugOrId := ctx.UserValue("slug_or_id").(string)
	if threadSlugOrId == "" {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	updatedThread, err := h.ThreadUCase.UpdateThreadDetails(threadSlugOrId, threadUpdate)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, updatedThread, http.StatusOK)
	case errors.ErrThreadNotFound:
		http_utils.SetJSONResponse(ctx, errors.ErrThreadNotFound, http.StatusNotFound)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}

func (h *ThreadHandler) UpdateThreadVote(ctx *fasthttp.RequestCtx) {
	threadVote := &models.ThreadVote{}
	if err := json.Unmarshal(ctx.PostBody(), threadVote); err != nil {
		http_utils.SetJSONResponse(ctx, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	threadSlugOrId := ctx.UserValue("slug_or_id").(string)
	if threadSlugOrId == "" {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	updatedThread, err := h.ThreadUCase.UpdateThreadVote(threadSlugOrId, threadVote)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, updatedThread, http.StatusOK)
	case errors.ErrThreadNotFound:
		http_utils.SetJSONResponse(ctx, errors.ErrThreadNotFound, http.StatusNotFound)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}
