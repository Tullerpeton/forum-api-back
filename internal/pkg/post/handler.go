package post

import "github.com/valyala/fasthttp"

type Handler interface {
	CreateNewPosts(ctx *fasthttp.RequestCtx)
	GetPostDetails(ctx *fasthttp.RequestCtx)
	GetPostsByThread(ctx *fasthttp.RequestCtx)
	UpdatePostDetails(ctx *fasthttp.RequestCtx)
}
