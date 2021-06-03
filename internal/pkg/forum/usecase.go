package forum

import "github.com/forum-api-back/internal/pkg/models"

type UseCase interface {
	CreateNewForum(forumInfo *models.ForumCreate) (*models.Forum, error)
	GetForumDetails(slug string) (*models.Forum, error)
}
