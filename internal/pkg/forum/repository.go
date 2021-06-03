package forum

import "github.com/forum-api-back/internal/pkg/models"

type Repository interface {
	InsertForum(forumInfo *models.ForumCreate) error
	SelectForumBySlug(forumSlug string) (*models.Forum, error)
}
