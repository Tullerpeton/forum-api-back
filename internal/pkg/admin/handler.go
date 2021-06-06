package admin

import "github.com/valyala/fasthttp"

type Handler interface {
	ClearBase(ctx *fasthttp.RequestCtx)
	GetBaseDetails(ctx *fasthttp.RequestCtx)
}
