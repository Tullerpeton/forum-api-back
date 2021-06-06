package usecase

import (
	"strconv"

	"github.com/forum-api-back/internal/pkg/forum"
	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/post"
	"github.com/forum-api-back/internal/pkg/thread"
	"github.com/forum-api-back/internal/pkg/user"
	"github.com/forum-api-back/pkg/errors"
)

type PostUseCase struct {
	PostRepo   post.Repository
	ThreadRepo thread.Repository
	ForumRepo  forum.Repository
	UserRepo   user.Repository
}

func NewUseCase(postRepo post.Repository, threadRepo thread.Repository,
	forumRepo forum.Repository, userRepo user.Repository) post.UseCase {
	return &PostUseCase{
		PostRepo:   postRepo,
		ThreadRepo: threadRepo,
		ForumRepo:  forumRepo,
		UserRepo:   userRepo,
	}
}

func (u *PostUseCase) CreateNewPosts(threadSlugOrId string, posts []*models.PostCreate) ([]*models.Post, error) {
	threadId, err := strconv.Atoi(threadSlugOrId)

	var newPosts []*models.Post
	if err != nil {
		selectedThread, errr := u.ThreadRepo.SelectThreadBySlug(threadSlugOrId)
		if errr != nil {
			return nil, errors.ErrThreadNotFound
		}
		if len(posts) == 0 {
			return []*models.Post{}, nil
		}
		newPosts, err = u.PostRepo.CreateNewPostsById(selectedThread.Id, selectedThread.Forum, posts)
	} else if threadId >= 1 {
		selectedThread, errr := u.ThreadRepo.SelectThreadById(uint64(threadId))
		if errr != nil {
			return nil, errors.ErrThreadNotFound
		}
		if len(posts) == 0 {
			return []*models.Post{}, nil
		}
		newPosts, err = u.PostRepo.CreateNewPostsById(uint64(threadId), selectedThread.Forum, posts)
	} else {
		return nil, errors.ErrPostNotFound
	}

	switch err {
	case nil:
		return newPosts, nil
	case errors.ErrUserNotFound:
		return nil, errors.ErrUserNotFound
	default:
		return nil, errors.ErrPostNotFound
	}
}

func (u *PostUseCase) GetPostDetail(postId uint64, related map[string]bool) (*models.PostDetails, error) {
	postDetails := &models.PostDetails{}

	selectedPost, err := u.PostRepo.SelectPostById(postId)
	if err != nil {
		return nil, errors.ErrThreadNotFound
	}
	postDetails.Post = selectedPost

	if related["user"] {
		selectedUser, err := u.UserRepo.SelectUserByNickName(selectedPost.Author)
		if err != nil {
			return nil, errors.ErrThreadNotFound
		}
		postDetails.Author = selectedUser
	}

	if related["thread"] {
		selectedThread, err := u.ThreadRepo.SelectThreadById(selectedPost.Thread)
		if err != nil {
			return nil, errors.ErrThreadNotFound
		}
		postDetails.Thread = selectedThread
	}

	if related["forum"] {
		selectedForum, err := u.ForumRepo.SelectForumBySlug(selectedPost.Forum)
		if err != nil {
			return nil, errors.ErrThreadNotFound
		}
		postDetails.Forum = selectedForum
	}

	return postDetails, nil
}

func (u *PostUseCase) GetPostsByThread(threadSlugOrId string, paginator *models.PostPaginator) ([]*models.Post, error) {
	threadId, err := strconv.Atoi(threadSlugOrId)

	var newPosts []*models.Post
	if err != nil {
		selectedThread, errr := u.ThreadRepo.SelectThreadBySlug(threadSlugOrId)
		if errr != nil {
			return nil, errors.ErrThreadNotFound
		}
		newPosts, err = u.PostRepo.SelectPostsById(selectedThread.Id, paginator)
	} else if threadId >= 1 {
		if _, errr := u.ThreadRepo.SelectThreadById(uint64(threadId)); errr != nil {
			return nil, errors.ErrThreadNotFound
		}
		newPosts, err = u.PostRepo.SelectPostsById(uint64(threadId), paginator)
	} else {
		return nil, errors.ErrThreadNotFound
	}

	if err != nil {
		return nil, errors.ErrThreadNotFound
	}

	return newPosts, nil
}

func (u *PostUseCase) UpdatePostDetails(postId uint64, postInfo *models.PostUpdate) (*models.Post, error) {
	selectedPost, err := u.PostRepo.SelectPostById(postId)
	if err != nil {
		return nil, errors.ErrPostNotFound
	}

	if postInfo.Message == "" || postInfo.Message == selectedPost.Message {
		return selectedPost, nil
	}

	err = u.PostRepo.UpdatePostById(postId, postInfo)
	if err != nil {
		return nil, errors.ErrPostNotFound
	}
	selectedPost.Message = postInfo.Message
	selectedPost.IsEdited = true

	return selectedPost, nil
}
