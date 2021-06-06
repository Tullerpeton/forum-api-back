package handler

import (
	"net/http"

	"github.com/forum-api-back/internal/pkg/admin"
	"github.com/forum-api-back/pkg/errors"
	"github.com/forum-api-back/pkg/tools/http_utils"

	"github.com/valyala/fasthttp"
)

type AdminHandler struct {
	AdminUCase admin.UseCase
}

func NewHandler(adminUCase admin.UseCase) admin.Handler {
	return &AdminHandler{
		AdminUCase: adminUCase,
	}
}

func (h *AdminHandler) ClearBase(ctx *fasthttp.RequestCtx) {
	err := h.AdminUCase.ClearBase()
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, struct{}{}, http.StatusOK)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}

func (h *AdminHandler) GetBaseDetails(ctx *fasthttp.RequestCtx) {
	baseDetails, err := h.AdminUCase.GetBaseDetails()
	switch err {
	case nil:
		http_utils.SetJSONResponse(ctx, baseDetails, http.StatusOK)
	default:
		http_utils.SetJSONResponse(ctx, errors.ErrInternalError, http.StatusInternalServerError)
	}
}
