package post

import "github.com/forum-api-back/internal/pkg/models"

type Repository interface {
	CreateNewPostsById(threadId uint64, forumSlug string,
		posts []*models.PostCreate) ([]*models.Post, error)
	SelectPostById(postId uint64) (*models.Post, error)
	SelectPostsById(threadId uint64, paginator *models.PostPaginator) ([]*models.Post, error)
	UpdatePostById(postId uint64, postInfo *models.PostUpdate) error
}
