package repository

import (
	"database/sql"

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
			"VALUES ($1, $2, $3, $4) "+
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
			"forum_slug, message, date_created "+
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
		&selectedThread.Author,
		&selectedThread.Forum,
		&selectedThread.Message,
		&selectedThread.DateCreated,
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
			"forum_slug, message, date_created "+
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
		&selectedThread.Author,
		&selectedThread.Forum,
		&selectedThread.Message,
		&selectedThread.DateCreated,
	)
	selectedThread.Slug = slug.String

	if err != nil {
		return nil, errors.ErrThreadNotFound
	}

	return selectedThread, nil
}

func (r *PostgresqlRepository) SelectThreadsByForum(forumSlug string,
	threadPaginator *models.ThreadPaginator) ([]*models.Thread, error) {
	var orderSort string
	if threadPaginator.SortOrder {
		orderSort = " DESC "
	} else {
		orderSort = " ASC "
	}

	rows, err := r.db.Query(
		"SELECT id, slug, title, author_nickname, "+
			"forum_slug, message, date_created "+
			"FROM threads "+
			"WHERE (forum_slug = $1 "+
			"date_created <= $2) "+
			"ORDER BY date_created "+orderSort+
			"LIMIT $3",
		forumSlug,
		threadPaginator.Since,
		threadPaginator.Limit,
	)

	if err != nil {
		return nil, errors.ErrThreadNotFound
	}

	threads := make([]*models.Thread, 0)
	for rows.Next() {
		selectedThread := &models.Thread{}
		slug := sql.NullString{}
		err := rows.Scan(
			&selectedThread.Id,
			&slug,
			&selectedThread.Title,
			&selectedThread.Author,
			&selectedThread.Forum,
			&selectedThread.Message,
			&selectedThread.DateCreated,
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
	row := r.db.QueryRow(
		"UPDATE threads SET "+
			"title = $1, "+
			"Message = $2 "+
			"WHERE slug = $3 "+
			"RETURNING id, slug, title, author_nickname, "+
			"	forum_slug, message, date_created",
		threadInfo.Title,
		threadInfo.Message,
		threadSlug,
	)

	updatedThread := &models.Thread{}
	slug := sql.NullString{}
	err := row.Scan(
		&updatedThread.Id,
		&slug,
		&updatedThread.Title,
		&updatedThread.Author,
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
	row := r.db.QueryRow(
		"UPDATE threads SET "+
			"title = $1, "+
			"Message = $2 "+
			"WHERE id = $3 "+
			"RETURNING id, slug, title, author_nickname, "+
			"	forum_slug, message, date_created",
		threadInfo.Title,
		threadInfo.Message,
		threadId,
	)

	updatedThread := &models.Thread{}
	slug := sql.NullString{}
	err := row.Scan(
		&updatedThread.Id,
		&slug,
		&updatedThread.Title,
		&updatedThread.Author,
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
		"WITH thread_info AS ( " +
			"	SELECT thread_id" +
			"	FROM threads " +
			" 	WHERE slug = $3 " +
			") " +
			"INSERT INTO votes (vote, author_nickname, thread_id) " +
			"SELECT $1, $2, thread_info.id " +
			"FROM thread_info "+
			"ON CONFLICT (thread_id, author_nickname) " +
			"DO UPDATE SET " +
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
			"INSERT INTO votes (vote, author_nickname, thread_id) " +
			"VALUES ($1, $2, $3) "+
			"ON CONFLICT (thread_id, author_nickname) " +
			"DO UPDATE SET " +
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
