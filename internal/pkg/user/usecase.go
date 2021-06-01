package user

import "github.com/forum-api-back/internal/pkg/models"

type UseCase interface {
	CreateNewUser(userInfo *models.User) (*models.User, error)
}
