package thread

import "github.com/valyala/fasthttp"

type Handler interface {
	CreateNewThread(ctx *fasthttp.RequestCtx)
	GetThreadsByForum(ctx *fasthttp.RequestCtx)
	GetThreadDetails(ctx *fasthttp.RequestCtx)
	UpdateThreadDetails(ctx *fasthttp.RequestCtx)
	UpdateThreadVote(ctx *fasthttp.RequestCtx)
}
