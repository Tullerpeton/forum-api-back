package post

import "github.com/valyala/fasthttp"

type Handler interface {
	CreateNewPosts(ctx *fasthttp.RequestCtx)
}
