package usecase

import (
	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/user"
	"github.com/forum-api-back/pkg/errors"
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
	err := u.UserRepo.InsertUser(userInfo)
	switch err {
	case nil:
		return userInfo, nil
	case errors.ErrDataConflict:
		selectedUser, err := u.UserRepo.SelectUserByEmailOrNickname(userInfo.Email, userInfo.NickName)
		if err != nil {
			return nil, errors.ErrInternalError
		}
		return selectedUser, errors.ErrDataConflict
	default:
		return nil, errors.ErrInternalError
	}
}

func (u *UserUseCase) GetUserByNickName(userNickName string) (*models.User, error) {
	selectedUser, err := u.UserRepo.SelectUserByNickName(userNickName)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	return selectedUser, err
}

func (u *UserUseCase) SetUserProfile(userInfo *models.User) (*models.User, error) {
	err := u.UserRepo.UpdateUserProfile(userInfo)
	switch err {
	case nil:
		return userInfo, nil
	case errors.ErrNotFoundInDB:
		return nil, errors.ErrUserNotFound
	case errors.ErrDataConflict:
		return nil, errors.ErrDataConflict
	default:
		return nil, errors.ErrInternalError
	}
}
