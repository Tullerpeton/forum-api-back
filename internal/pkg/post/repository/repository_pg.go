package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/post"
	"github.com/forum-api-back/pkg/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) post.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

func (r *PostgresqlRepository) CreateNewPostsById(threadId uint64, forumSlug string,
	posts []*models.PostCreate) ([]*models.Post, error) {
	countPosts := len(posts)
	insertData := make([]string, countPosts)
	for i := 0; i < countPosts; i++ {
		insertData[i] = fmt.Sprintf(
			"(%d, '%s', '%s', '%s', %d)",
			posts[i].Parent,
			posts[i].Author,
			posts[i].Message,
			forumSlug,
			threadId,
		)
	}

	rows, err := r.db.Query(
		"INSERT INTO posts (parent_message_id, author_nickname, message, " +
			"forum_slug, thread_id) VALUES " +
			strings.Join(insertData, ", ") +
			" RETURNING id, date_created",
	)

	if err != nil {
		if strings.Contains(err.Error(), "posts_author_nickname") {
			return nil, errors.ErrUserNotFound
		}

		return nil, errors.ErrInternalError
	}
	defer rows.Close()

	newPosts := make([]*models.Post, countPosts, countPosts)
	for i := 0; i < countPosts && rows.Next(); i++ {
		newPosts[i] = &models.Post{}
		newPosts[i].Message = posts[i].Message
		newPosts[i].Author = posts[i].Author
		newPosts[i].Parent = posts[i].Parent
		newPosts[i].Thread = threadId
		newPosts[i].Forum = forumSlug

		err = rows.Scan(
			&newPosts[i].Id,
			&newPosts[i].DateCreated,
		)

		if err != nil {
			return nil, errors.ErrInternalError
		}

	}

	return newPosts, nil
}

func (r *PostgresqlRepository) SelectPostById(postId uint64) (*models.Post, error) {
	row := r.db.QueryRow(
		"SELECT id, parent_message_id, author_nickname, message, "+
			"is_edited, forum_slug, thread_id, date_created "+
			"FROM posts "+
			"WHERE id = $1",
		postId,
	)

	selectedPost := &models.Post{}
	err := row.Scan(
		&selectedPost.Id,
		&selectedPost.Parent,
		&selectedPost.Author,
		&selectedPost.Message,
		&selectedPost.IsEdited,
		&selectedPost.Forum,
		&selectedPost.Thread,
		&selectedPost.DateCreated,
	)

	if err != nil {
		return nil, errors.ErrPostNotFound
	}

	return selectedPost, nil
}

func (r *PostgresqlRepository) SelectPostsById(threadId uint64, paginator *models.PostPaginator) ([]*models.Post, error) {
	var orderSort, orderCompare string
	if paginator.SortOrder {
		orderSort = " DESC "
		orderCompare = " < "
	} else {
		orderSort = " ASC "
		orderCompare = " > "
	}

	var rows *sql.Rows
	var err error
	switch paginator.Sort {
	case "flat":
		if paginator.Since != 0 {
			rows, err = r.db.Query(
				"SELECT id, parent_message_id, author_nickname, message, "+
					"is_edited, forum_slug, thread_id, date_created "+
					"FROM posts "+
					"WHERE (thread_id = $1 AND id"+orderCompare+"$2) "+
					"ORDER BY id "+orderSort+
					"LIMIT $3",
				threadId,
				paginator.Since,
				paginator.Limit,
			)
		} else {
			rows, err = r.db.Query(
				"SELECT id, parent_message_id, author_nickname, message, "+
					"is_edited, forum_slug, thread_id, date_created "+
					"FROM posts "+
					"WHERE thread_id = $1"+
					"ORDER BY id "+orderSort+
					"LIMIT $2",
				threadId,
				paginator.Limit,
			)
		}
	case "tree":
		if paginator.Since != 0 {
			rows, err = r.db.Query(
				"SELECT p1.id, p1.parent_message_id, p1.author_nickname, p1.message, "+
					"p1.is_edited, p1.forum_slug, p1.thread_id, p1.date_created "+
					"FROM posts p1 "+
					"JOIN posts p2 ON (p2.id = $2) "+
					"WHERE (p1.thread_id = $1 AND p1.path_of_nesting"+orderCompare+"p2.path_of_nesting) "+
					"ORDER BY p1.path_of_nesting[1]"+orderSort+", "+
					"	p1.path_of_nesting "+orderSort+
					"LIMIT $3",
				threadId,
				paginator.Since,
				paginator.Limit,
			)
		} else {
			rows, err = r.db.Query(
				"SELECT p1.id, p1.parent_message_id, p1.author_nickname, p1.message, "+
					"p1.is_edited, p1.forum_slug, p1.thread_id, p1.date_created "+
					"FROM posts p1 "+
					"WHERE (p1.thread_id = $1) "+
					"ORDER BY p1.path_of_nesting[1]"+orderSort+", "+
					"	p1.path_of_nesting "+orderSort+
					"LIMIT $2",
				threadId,
				paginator.Limit,
			)
		}
	case "parent_tree":
		if paginator.Since != 0 {
			rows, err = r.db.Query(
				"SELECT p1.id, p1.parent_message_id, p1.author_nickname, p1.message, "+
					"p1.is_edited, p1.forum_slug, p1.thread_id, p1.date_created "+
					"FROM posts p1 "+
					"WHERE p1.path_of_nesting[1] IN ( "+
					"	SELECT id "+
					"	FROM posts "+
					"	WHERE (thread_id = $1 AND parent_message_id = 0 AND "+
					"		path_of_nesting[1] "+orderCompare+" ("+
					"			SELECT path_of_nesting[1] "+
					"			FROM posts "+
					"			WHERE id = $2 "+
					"		) "+
					"	) "+
					"	ORDER BY id "+orderSort+
					"	LIMIT $3"+
					") "+
					"ORDER BY p1.path_of_nesting[1] "+orderSort+", p1.path_of_nesting",
				threadId,
				paginator.Since,
				paginator.Limit,
			)
		} else {
			rows, err = r.db.Query(
				"SELECT p1.id, p1.parent_message_id, p1.author_nickname, p1.message, "+
					"p1.is_edited, p1.forum_slug, p1.thread_id, p1.date_created "+
					"FROM posts p1 "+
					"WHERE p1.path_of_nesting[1] IN ( "+
					"	SELECT id "+
					"	FROM posts "+
					"	WHERE (thread_id = $1 AND parent_message_id = 0) "+
					"	ORDER BY id "+orderSort+
					"	LIMIT $2"+
					" ) "+
					"ORDER BY p1.path_of_nesting[1] "+orderSort+", p1.path_of_nesting",
				threadId,
				paginator.Limit,
			)
		}
	default:
		return nil, errors.ErrPostNotFound
	}

	if err != nil {
		return nil, errors.ErrPostNotFound
	}
	defer rows.Close()

	posts := make([]*models.Post, 0)
	for rows.Next() {
		selectedPost := &models.Post{}
		err = rows.Scan(
			&selectedPost.Id,
			&selectedPost.Parent,
			&selectedPost.Author,
			&selectedPost.Message,
			&selectedPost.IsEdited,
			&selectedPost.Forum,
			&selectedPost.Thread,
			&selectedPost.DateCreated,
		)
		if err != nil {
			return nil, errors.ErrPostNotFound
		}

		posts = append(posts, selectedPost)
	}

	return posts, nil
}

func (r *PostgresqlRepository) UpdatePostById(postId uint64, postInfo *models.PostUpdate) error {
	if postInfo.Message == "" {
		return nil
	}

	_, err := r.db.Exec(
		"UPDATE posts SET "+
			"message = $1, "+
			"is_edited = true "+
			"WHERE id = $2",
		postInfo.Message,
		postId,
	)

	if err != nil {
		return errors.ErrDataConflict
	}

	return nil
}
