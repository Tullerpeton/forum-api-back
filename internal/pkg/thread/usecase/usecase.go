package usecase

import (
	"strconv"

	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/thread"
	"github.com/forum-api-back/pkg/errors"
)

type ThreadUseCase struct {
	ThreadRepo thread.Repository
}

func NewUseCase(threadRepo thread.Repository) thread.UseCase {
	return &ThreadUseCase{
		ThreadRepo: threadRepo,
	}
}

func (u *ThreadUseCase) CreateNewThread(forumSlug string,
	threadInfo *models.ThreadCreate) (*models.Thread, error) {
	threadId, err := u.ThreadRepo.InsertThread(forumSlug, threadInfo)
	switch err {
	case nil:
		return &models.Thread{
			Id:          threadId,
			Title:       threadInfo.Title,
			Author:      threadInfo.AuthorNickName,
			Forum:       forumSlug,
			Message:     threadInfo.Message,
			Slug:        threadInfo.Slug,
			DateCreated: threadInfo.DateCreated,
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
	threads, err := u.ThreadRepo.SelectThreadsByForum(forumSlug, threadPaginator)
	if err != nil {
		return nil, errors.ErrForumNotFound
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
	} else if threadId >= 1 {
		updatedThread, err = u.ThreadRepo.UpdateThreadDetailsById(uint64(threadId), threadInfo)
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
