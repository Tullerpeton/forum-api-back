package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/post"
	"github.com/forum-api-back/pkg/errors"
	"github.com/forum-api-back/pkg/tools/http_utils"

	"github.com/valyala/fasthttp"
)

type PostHandler struct {
	PostUCase post.UseCase
}

func NewHandler(postUCase post.UseCase) post.Handler {
	return &PostHandler{
		PostUCase: postUCase,
	}
}

func (h *PostHandler) CreateNewPosts(ctx *fasthttp.RequestCtx) {
	var postsInfo []*models.PostCreate
	if err := json.Unmarshal(ctx.PostBody(), &postsInfo); err != nil {
		http_utils.SetJSONResponse(ctx, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	threadSlugOrId := ctx.UserValue("slug_or_id").(string)
	if threadSlugOrId == "" {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	newThreads, err := h.PostUCase.CreateNewPosts(threadSlugOrId, postsInfo)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, newThreads, http.StatusCreated)
	case errors.ErrPostNotFound:
		http_utils.SetJSONResponse(ctx, errors.ErrDataConflict, http.StatusConflict)
	case errors.ErrUserNotFound:
		http_utils.SetJSONResponse(ctx, errors.ErrUserNotFound, http.StatusNotFound)
	case errors.ErrThreadNotFound:
		http_utils.SetJSONResponse(ctx,
			errors.Error{Message: fmt.Sprintf("Can't find post thread by id: %s", threadSlugOrId)}, http.StatusNotFound)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetPostDetails(ctx *fasthttp.RequestCtx) {
	forumSlug, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil || forumSlug < 1 {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	related := map[string]bool{"user": false, "forum": false, "thread": false}
	queryRow := ctx.QueryArgs().Peek("related")
	for key, _ := range related {
		if bytes.Contains(queryRow, []byte(key)) {
			related[key] = true
		}
	}

	postDetails, err := h.PostUCase.GetPostDetail(uint64(forumSlug), related)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, postDetails, http.StatusOK)
	case errors.ErrThreadNotFound:
		http_utils.SetJSONResponse(ctx, errors.ErrThreadNotFound, http.StatusNotFound)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}

func (h *PostHandler) GetPostsByThread(ctx *fasthttp.RequestCtx) {
	threadSlugOrId := ctx.UserValue("slug_or_id").(string)
	if threadSlugOrId == "" {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	postPaginator := &models.PostPaginator{Limit: 100, Sort: "flat"}
	parseTime, err := strconv.Atoi(string(ctx.FormValue("since")))
	if err == nil {
		postPaginator.Since = uint64(parseTime)
	}

	if isDesc := string(ctx.FormValue("desc")); isDesc == "true" {
		postPaginator.SortOrder = true
	}

	parseLimit, err := strconv.Atoi(string(ctx.FormValue("limit")))
	if err == nil {
		postPaginator.Limit = uint64(parseLimit)
	}

	if sort := string(ctx.FormValue("sort")); sort != "" {
		postPaginator.Sort = sort
	}

	selectedPosts, err := h.PostUCase.GetPostsByThread(threadSlugOrId, postPaginator)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, selectedPosts, http.StatusOK)
	case errors.ErrThreadNotFound:
		http_utils.SetJSONResponse(ctx, errors.ErrThreadNotFound, http.StatusNotFound)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}

func (h *PostHandler) UpdatePostDetails(ctx *fasthttp.RequestCtx) {
	var postsInfo *models.PostUpdate
	if err := json.Unmarshal(ctx.PostBody(), &postsInfo); err != nil {
		http_utils.SetJSONResponse(ctx, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	forumId, err := strconv.Atoi(ctx.UserValue("id").(string))
	if err != nil || forumId < 1 {
		http_utils.SetJSONResponse(ctx, errors.ErrBadArguments, http.StatusBadRequest)
		return
	}

	updatedPost, err := h.PostUCase.UpdatePostDetails(uint64(forumId), postsInfo)
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, updatedPost, http.StatusOK)
	case errors.ErrPostNotFound:
		http_utils.SetJSONResponse(ctx, errors.ErrPostNotFound, http.StatusNotFound)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}
