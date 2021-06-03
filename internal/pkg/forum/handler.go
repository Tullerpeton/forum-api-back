package forum

import (
	"github.com/valyala/fasthttp"
)

type Handler interface {
	CreateNewForum(ctx *fasthttp.RequestCtx)
	GetForumDetails(ctx *fasthttp.RequestCtx)
}
