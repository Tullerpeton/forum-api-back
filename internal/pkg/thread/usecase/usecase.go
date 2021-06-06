package usecase

import (
	"strconv"

	"github.com/forum-api-back/internal/pkg/forum"
	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/thread"
	"github.com/forum-api-back/pkg/errors"
)

type ThreadUseCase struct {
	ThreadRepo thread.Repository
	ForumRepo  forum.Repository
}

func NewUseCase(threadRepo thread.Repository, forumRepo forum.Repository) thread.UseCase {
	return &ThreadUseCase{
		ThreadRepo: threadRepo,
		ForumRepo:  forumRepo,
	}
}

func (u *ThreadUseCase) CreateNewThread(forumSlug string,
	threadInfo *models.ThreadCreate) (*models.Thread, error) {
	selectedForum, err := u.ForumRepo.SelectForumBySlug(forumSlug)
	if err != nil {
		return nil, errors.ErrDataConflict
	}

	threadId, err := u.ThreadRepo.InsertThread(selectedForum.Slug, threadInfo)
	switch err {
	case nil:
		return &models.Thread{
			Id:             threadId,
			Title:          threadInfo.Title,
			AuthorNickName: threadInfo.AuthorNickName,
			Forum:          selectedForum.Slug,
			Message:        threadInfo.Message,
			Slug:           threadInfo.Slug,
			DateCreated:    threadInfo.DateCreated,
		}, nil
	case errors.ErrDataConflict:
		if threadInfo.Slug != "" {
			selectedThread, err := u.ThreadRepo.SelectThreadBySlug(threadInfo.Slug)
			if err == nil {
				return selectedThread, errors.ErrAlreadyExists
			}
		}
		return nil, errors.ErrDataConflict
	default:
		return nil, errors.ErrInternalError
	}
}

func (u *ThreadUseCase) GetThreadsByForum(forumSlug string,
	threadPaginator *models.ThreadPaginator) ([]*models.Thread, error) {
	if _, err := u.ForumRepo.SelectForumBySlug(forumSlug); err != nil {
		return nil, errors.ErrForumNotFound
	}

	threads, err := u.ThreadRepo.SelectThreadsByForum(forumSlug, threadPaginator)
	if err != nil {
		return nil, errors.ErrInternalError
	}

	return threads, nil
}

func (u *ThreadUseCase) GetThreadDetails(threadSlugOrId string) (*models.Thread, error) {
	threadId, err := strconv.Atoi(threadSlugOrId)

	var selectedThread *models.Thread
	if err != nil {
		selectedThread, err = u.ThreadRepo.SelectThreadBySlug(threadSlugOrId)
	} else if threadId >= 1 {
		selectedThread, err = u.ThreadRepo.SelectThreadById(uint64(threadId))
	} else {
		return nil, errors.ErrThreadNotFound
	}

	if err != nil {
		return nil, errors.ErrThreadNotFound
	}

	return selectedThread, nil
}

func (u *ThreadUseCase) UpdateThreadDetails(threadSlugOrId string,
	threadInfo *models.ThreadUpdate) (*models.Thread, error) {
	threadId, err := strconv.Atoi(threadSlugOrId)

	var updatedThread *models.Thread
	if err != nil {
		updatedThread, err = u.ThreadRepo.UpdateThreadDetailsBySlug(threadSlugOrId, threadInfo)
		if err == errors.ErrEmptyParameters {
			updatedThread, err = u.ThreadRepo.SelectThreadBySlug(threadSlugOrId)
		}
	} else if threadId >= 1 {
		updatedThread, err = u.ThreadRepo.UpdateThreadDetailsById(uint64(threadId), threadInfo)
		if err == errors.ErrEmptyParameters {
			updatedThread, err = u.ThreadRepo.SelectThreadById(uint64(threadId))
		}
	} else {
		return nil, errors.ErrThreadNotFound
	}

	if err != nil {
		return nil, errors.ErrThreadNotFound
	}
	return updatedThread, nil
}

func (u *ThreadUseCase) UpdateThreadVote(threadSlugOrId string,
	threadVote *models.ThreadVote) (*models.Thread, error) {
	threadId, err := strconv.Atoi(threadSlugOrId)

	var updatedThread *models.Thread
	if err != nil {
		if err = u.ThreadRepo.UpdateThreadVoteBySlug(threadSlugOrId, threadVote); err != nil {
			return nil, errors.ErrThreadNotFound
		}
		updatedThread, err = u.ThreadRepo.SelectThreadBySlug(threadSlugOrId)
	} else if threadId >= 1 {
		if err = u.ThreadRepo.UpdateThreadVoteById(uint64(threadId), threadVote); err != nil {
			return nil, errors.ErrThreadNotFound
		}
		updatedThread, err = u.ThreadRepo.SelectThreadById(uint64(threadId))
	} else {
		return nil, errors.ErrThreadNotFound
	}

	if err != nil {
		return nil, errors.ErrThreadNotFound
	}

	return updatedThread, nil
}
