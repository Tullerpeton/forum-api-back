package usecase

import (
	"github.com/forum-api-back/internal/pkg/forum"
	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/user"
	"github.com/forum-api-back/pkg/errors"
)

type UserUseCase struct {
	UserRepo  user.Repository
	ForumRepo forum.Repository
}

func NewUseCase(userRepo user.Repository, forumRepo forum.Repository) user.UseCase {
	return &UserUseCase{
		UserRepo:  userRepo,
		ForumRepo: forumRepo,
	}
}

func (u *UserUseCase) CreateNewUser(userInfo *models.User) ([]*models.User, error) {
	err := u.UserRepo.InsertUser(userInfo)
	switch err {
	case nil:
		return []*models.User{userInfo}, nil
	case errors.ErrDataConflict:
		selectedUser, err := u.UserRepo.SelectUserByEmailOrNickname(userInfo.Email, userInfo.NickName)
		if err != nil {
			return nil, errors.ErrInternalError
		}
		return selectedUser, errors.ErrAlreadyExists
	default:
		return nil, errors.ErrInternalError
	}
}

func (u *UserUseCase) GetUserByNickName(userNickName string) (*models.User, error) {
	selectedUser, err := u.UserRepo.SelectUserByNickName(userNickName)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	return selectedUser, nil
}

func (u *UserUseCase) GetUsersByForum(forumSlug string, paginator *models.UserPaginator) ([]*models.User, error) {
	if _, err := u.ForumRepo.SelectForumBySlug(forumSlug); err != nil {
		return nil, errors.ErrForumNotFound
	}

	selectedUsers, err := u.UserRepo.SelectUsersByForum(forumSlug, paginator)
	switch err {
	case nil:
		return selectedUsers, nil
	case errors.ErrNotFoundInDB:
		return nil, errors.ErrUserNotFound
	default:
		return nil, errors.ErrInternalError
	}
}

func (u *UserUseCase) SetUserProfile(userInfo *models.User) (*models.User, error) {
	selectedUser, err := u.UserRepo.SelectUserByNickName(userInfo.NickName)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	err = u.UserRepo.UpdateUserProfile(userInfo)
	switch err {
	case nil:
		if userInfo.Email != "" {
			selectedUser.Email = userInfo.Email
		}
		if userInfo.About != "" {
			selectedUser.About = userInfo.About
		}
		if userInfo.FullName != "" {
			selectedUser.FullName = userInfo.FullName
		}
		return selectedUser, nil
	case errors.ErrNotFoundInDB:
		return nil, errors.ErrUserNotFound
	case errors.ErrDataConflict:
		return nil, errors.ErrAlreadyExists
	default:
		return nil, errors.ErrInternalError
	}
}
