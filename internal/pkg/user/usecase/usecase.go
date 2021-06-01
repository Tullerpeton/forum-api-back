package usecase

import (
	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/user"
)

type UserUseCase struct {
	UserRepo user.Repository
}

func NewUseCase(userRepo user.Repository) user.UseCase {
	return &UserUseCase{
		UserRepo: userRepo,
	}
}

func (u *UserUseCase) CreateNewUser(userInfo *models.User) (*models.User, error) {

}
