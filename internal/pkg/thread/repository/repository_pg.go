package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/thread"
	"github.com/forum-api-back/pkg/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) thread.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

func (r *PostgresqlRepository) InsertThread(forumSlug string,
	threadInfo *models.ThreadCreate) (uint64, error) {
	slug := sql.NullString{}
	if threadInfo.Slug != "" {
		slug.String = threadInfo.Slug
		slug.Valid = true
	}
	row := r.db.QueryRow(
		"INSERT INTO threads(slug, title, author_nickname, "+
			"	forum_slug, message, date_created) "+
			"VALUES ($1, $2, $3, $4, $5, $6) "+
			"RETURNING id",
		slug,
		threadInfo.Title,
		threadInfo.AuthorNickName,
		forumSlug,
		threadInfo.Message,
		threadInfo.DateCreated,
	)

	var threadId uint64
	if err := row.Scan(&threadId); err != nil {
		return 0, errors.ErrDataConflict
	}

	return threadId, nil
}

func (r *PostgresqlRepository) SelectThreadBySlug(threadSlug string) (*models.Thread, error) {
	row := r.db.QueryRow(
		"SELECT id, slug, title, author_nickname, "+
			"forum_slug, message, date_created, votes "+
			"FROM threads "+
			"WHERE slug = $1",
		threadSlug,
	)

	selectedThread := &models.Thread{}
	slug := sql.NullString{}
	err := row.Scan(
		&selectedThread.Id,
		&slug,
		&selectedThread.Title,
		&selectedThread.AuthorNickName,
		&selectedThread.Forum,
		&selectedThread.Message,
		&selectedThread.DateCreated,
		&selectedThread.Votes,
	)
	selectedThread.Slug = slug.String

	if err != nil {
		return nil, errors.ErrThreadNotFound
	}

	return selectedThread, nil
}

func (r *PostgresqlRepository) SelectThreadById(threadId uint64) (*models.Thread, error) {
	row := r.db.QueryRow(
		"SELECT id, slug, title, author_nickname, "+
			"forum_slug, message, date_created, votes "+
			"FROM threads "+
			"WHERE id = $1",
		threadId,
	)

	selectedThread := &models.Thread{}
	slug := sql.NullString{}
	err := row.Scan(
		&selectedThread.Id,
		&slug,
		&selectedThread.Title,
		&selectedThread.AuthorNickName,
		&selectedThread.Forum,
		&selectedThread.Message,
		&selectedThread.DateCreated,
		&selectedThread.Votes,
	)
	selectedThread.Slug = slug.String

	if err != nil {
		return nil, errors.ErrThreadNotFound
	}

	return selectedThread, nil
}

func (r *PostgresqlRepository) SelectThreadsByForum(forumSlug string,
	threadPaginator *models.ThreadPaginator) ([]*models.Thread, error) {
	var orderSort, orderCompare string
	if threadPaginator.SortOrder {
		orderSort = " DESC "
		orderCompare = " <= "
	} else {
		orderSort = " ASC "
		orderCompare = " >= "
	}

	var rows *sql.Rows
	var err error
	if threadPaginator.Since.IsZero() {
		rows, err = r.db.Query(
			"SELECT id, slug, title, author_nickname, "+
				"forum_slug, message, date_created, votes "+
				"FROM threads "+
				"WHERE forum_slug = $1 "+
				"ORDER BY date_created "+orderSort+
				"LIMIT $2",
			forumSlug,
			threadPaginator.Limit,
		)
	} else {
		rows, err = r.db.Query(
			"SELECT id, slug, title, author_nickname, "+
				"forum_slug, message, date_created, votes "+
				"FROM threads "+
				"WHERE (forum_slug = $1 AND "+
				"date_created"+orderCompare+"$2) "+
				"ORDER BY date_created "+orderSort+
				"LIMIT $3",
			forumSlug,
			threadPaginator.Since,
			threadPaginator.Limit,
		)
	}

	if err != nil {
		return nil, errors.ErrThreadNotFound
	}
	defer rows.Close()

	threads := make([]*models.Thread, 0)
	for rows.Next() {
		selectedThread := &models.Thread{}
		slug := sql.NullString{}
		err := rows.Scan(
			&selectedThread.Id,
			&slug,
			&selectedThread.Title,
			&selectedThread.AuthorNickName,
			&selectedThread.Forum,
			&selectedThread.Message,
			&selectedThread.DateCreated,
			&selectedThread.Votes,
		)
		selectedThread.Slug = slug.String
		if err != nil {
			return nil, errors.ErrInternalError
		}

		threads = append(threads, selectedThread)
	}

	return threads, nil
}

func (r *PostgresqlRepository) UpdateThreadDetailsBySlug(threadSlug string,
	threadInfo *models.ThreadUpdate) (*models.Thread, error) {
	columns := make([]string, 0)
	args := make([]interface{}, 1)
	args[0] = threadSlug
	if threadInfo.Title != "" {
		args = append(args, threadInfo.Title)
		columns = append(columns, fmt.Sprintf("title = $%d ", len(args)))
	}
	if threadInfo.Message != "" {
		args = append(args, threadInfo.Message)
		columns = append(columns, fmt.Sprintf("message = $%d ", len(args)))
	}

	if len(args) == 1 {
		return nil, errors.ErrEmptyParameters
	}

	row := r.db.QueryRow(
		"UPDATE threads SET "+
			strings.Join(columns, ", ")+
			" WHERE slug = $1 "+
			"RETURNING id, slug, title, author_nickname, "+
			"	forum_slug, message, date_created",
		args...,
	)

	updatedThread := &models.Thread{}
	slug := sql.NullString{}
	err := row.Scan(
		&updatedThread.Id,
		&slug,
		&updatedThread.Title,
		&updatedThread.AuthorNickName,
		&updatedThread.Forum,
		&updatedThread.Message,
		&updatedThread.DateCreated,
	)
	updatedThread.Slug = slug.String

	if err != nil {
		return nil, errors.ErrThreadNotFound
	}

	return updatedThread, nil
}

func (r *PostgresqlRepository) UpdateThreadDetailsById(threadId uint64,
	threadInfo *models.ThreadUpdate) (*models.Thread, error) {

	columns := make([]string, 0)
	args := make([]interface{}, 1)
	args[0] = threadId
	if threadInfo.Title != "" {
		args = append(args, threadInfo.Title)
		columns = append(columns, fmt.Sprintf("title = $%d ", len(args)))
	}
	if threadInfo.Message != "" {
		args = append(args, threadInfo.Message)
		columns = append(columns, fmt.Sprintf("message = $%d ", len(args)))
	}

	if len(args) == 1 {
		return nil, errors.ErrEmptyParameters
	}

	row := r.db.QueryRow(
		"UPDATE threads SET "+
			strings.Join(columns, ", ")+
			" WHERE id = $1 "+
			"RETURNING id, slug, title, author_nickname, "+
			"	forum_slug, message, date_created",
		args...,
	)

	updatedThread := &models.Thread{}
	slug := sql.NullString{}
	err := row.Scan(
		&updatedThread.Id,
		&slug,
		&updatedThread.Title,
		&updatedThread.AuthorNickName,
		&updatedThread.Forum,
		&updatedThread.Message,
		&updatedThread.DateCreated,
	)
	updatedThread.Slug = slug.String

	if err != nil {
		return nil, errors.ErrThreadNotFound
	}

	return updatedThread, nil
}

func (r *PostgresqlRepository) UpdateThreadVoteBySlug(threadSlug string,
	threadVote *models.ThreadVote) error {
	_, err := r.db.Exec(
		"WITH thread_info AS ( "+
			"	SELECT id "+
			"	FROM threads "+
			" 	WHERE slug = $3 "+
			") "+
			"INSERT INTO votes (vote, author_nickname, thread_id) "+
			"SELECT $1, $2, thread_info.id "+
			"FROM thread_info "+
			"ON CONFLICT (thread_id, author_nickname) "+
			"DO UPDATE SET "+
			"vote = $1",
		threadVote.Voice,
		threadVote.NickName,
		threadSlug,
	)

	if err != nil {
		return errors.ErrDataConflict
	}

	return nil
}

func (r *PostgresqlRepository) UpdateThreadVoteById(threadId uint64,
	threadVote *models.ThreadVote) error {
	_, err := r.db.Exec(
		"INSERT INTO votes (vote, author_nickname, thread_id) "+
			"VALUES ($1, $2, $3) "+
			"ON CONFLICT (thread_id, author_nickname) "+
			"DO UPDATE SET "+
			"vote = $1",
		threadVote.Voice,
		threadVote.NickName,
		threadId,
	)

	if err != nil {
		return errors.ErrDataConflict
	}

	return nil
}
