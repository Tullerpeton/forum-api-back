package thread

import "github.com/forum-api-back/internal/pkg/models"

type Repository interface {
	InsertThread(forumSlug string,
		threadInfo *models.ThreadCreate) (uint64, error)
	SelectThreadBySlug(threadSlug string) (*models.Thread, error)
	SelectThreadById(threadId uint64) (*models.Thread, error)
	SelectThreadsByForum(forumSlug string,
		threadPaginator *models.ThreadPaginator) ([]*models.Thread, error)
	UpdateThreadDetailsBySlug(threadSlug string,
		threadInfo *models.ThreadUpdate) (*models.Thread, error)
	UpdateThreadDetailsById(threadId uint64,
		threadInfo *models.ThreadUpdate) (*models.Thread, error)
	UpdateThreadVoteBySlug(threadSlug string,
		threadVote *models.ThreadVote) (*models.Thread, error)
	UpdateThreadVoteById(threadId uint64,
		threadVote *models.ThreadVote) (*models.Thread, error)
}
