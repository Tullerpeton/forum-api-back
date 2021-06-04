package thread

import "github.com/forum-api-back/internal/pkg/models"

type UseCase interface {
	CreateNewThread(forumSlug string,
		threadInfo *models.ThreadCreate) (*models.Thread, error)
	GetThreadsByForum(forumSlug string,
		threadPaginator *models.ThreadPaginator) ([]*models.Thread, error)
	GetThreadDetails(threadSlugOrId string) (*models.Thread, error)
	UpdateThreadDetails(threadSlugOrId string,
		threadInfo *models.ThreadUpdate) (*models.Thread, error)
	UpdateThreadVote(threadSlugOrId string,
		threadVote *models.ThreadVote) (*models.Thread, error)
}
