package handler

import "github.com/forum-api-back/internal/pkg/post"

type PostHandler struct {
	PostUCase post.UseCase
}

func NewHandler(postUCase post.UseCase) post.Handler {
	return &PostHandler{
		PostUCase: postUCase,
	}
}


