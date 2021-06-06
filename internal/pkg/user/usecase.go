package user

import "github.com/forum-api-back/internal/pkg/models"

type UseCase interface {
	CreateNewUser(userInfo *models.User) ([]*models.User, error)
	GetUserByNickName(userNickName string) (*models.User, error)
	GetUsersByForum(forumSlug string, paginator *models.UserPaginator) ([]*models.User, error)
	SetUserProfile(userInfo *models.User) (*models.User, error)
}
