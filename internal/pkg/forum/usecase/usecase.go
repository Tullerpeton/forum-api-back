package usecase

import (
	"github.com/forum-api-back/internal/pkg/forum"
	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/user"
	"github.com/forum-api-back/pkg/errors"
)

type ForumUseCase struct {
	ForumRepo forum.Repository
	UserRepo  user.Repository
}

func NewUseCase(forumRepo forum.Repository, userRepo user.Repository) forum.UseCase {
	return &ForumUseCase{
		ForumRepo: forumRepo,
		UserRepo:  userRepo,
	}
}

func (u *ForumUseCase) CreateNewForum(forumInfo *models.ForumCreate) (*models.Forum, error) {
	author, err := u.UserRepo.SelectUserByNickName(forumInfo.AuthorNickName)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	forumInfo.AuthorNickName = author.NickName

	err = u.ForumRepo.InsertForum(forumInfo)
	switch err {
	case nil:
		return &models.Forum{
			Title:          forumInfo.Title,
			AuthorNickName: forumInfo.AuthorNickName,
			Slug:           forumInfo.Slug,
		}, nil
	case errors.ErrDataConflict:
		selectedForum, err := u.ForumRepo.SelectForumBySlug(forumInfo.Slug)
		if err != nil {
			return nil, errors.ErrInternalError
		}
		return selectedForum, errors.ErrAlreadyExists
	default:
		return nil, errors.ErrInternalError
	}
}

func (u *ForumUseCase) GetForumDetails(forumSlug string) (*models.Forum, error) {
	selectedForum, err := u.ForumRepo.SelectForumBySlug(forumSlug)
	if err != nil {
		return nil, errors.ErrForumNotFound
	}

	return selectedForum, nil
}
