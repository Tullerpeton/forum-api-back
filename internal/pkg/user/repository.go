package user

import "github.com/forum-api-back/internal/pkg/models"

type Repository interface {
	InsertUser(userInfo *models.User) error
	SelectUserByEmailOrNickname(email, nickname string) (*models.User, error)
	SelectUserByNickName(nickname string) (*models.User, error)
	UpdateUserProfile(userInfo *models.User) error
}
