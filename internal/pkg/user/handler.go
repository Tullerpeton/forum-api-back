package user

import (
	"github.com/valyala/fasthttp"
)

type Handler interface {
	CreateNewUser(ctx *fasthttp.RequestCtx)
	GetUserProfile(ctx *fasthttp.RequestCtx)
	GetUsersByForum(ctx *fasthttp.RequestCtx)
	UpdateUserProfile(ctx *fasthttp.RequestCtx)
}
