package post

import "github.com/forum-api-back/internal/pkg/models"

type UseCase interface {
	CreateNewPosts(threadSlugOrId string, posts []*models.PostCreate) ([]*models.Post, error)
	GetPostDetail(postId uint64, related map[string]bool) (*models.PostDetails, error)
	GetPostsByThread(threadSlugOrId string, paginator *models.PostPaginator) ([]*models.Post, error)
	UpdatePostDetails(postId uint64, postInfo *models.PostUpdate) (*models.Post, error)
}
