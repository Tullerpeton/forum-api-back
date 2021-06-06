package repository

import (
	"database/sql"

	"github.com/forum-api-back/internal/pkg/forum"
	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/pkg/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) forum.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

func (r *PostgresqlRepository) InsertForum(forumInfo *models.ForumCreate) error {
	_, err := r.db.Exec(
		"INSERT INTO forums(title, author_nickname, slug) "+
			"VALUES ($1, $2, $3)",
		forumInfo.Title,
		forumInfo.AuthorNickName,
		forumInfo.Slug,
	)

	if err != nil {
		return errors.ErrDataConflict
	}

	return nil
}

func (r *PostgresqlRepository) SelectForumBySlug(forumSlug string) (*models.Forum, error) {
	row := r.db.QueryRow(
		"SELECT title, author_nickname, slug, count_posts, count_threads "+
			"FROM forums "+
			"WHERE slug = $1",
		forumSlug,
	)

	selectedForum := &models.Forum{}
	err := row.Scan(
		&selectedForum.Title,
		&selectedForum.AuthorNickName,
		&selectedForum.Slug,
		&selectedForum.Posts,
		&selectedForum.Threads,
	)

	switch err {
	case nil:
		return selectedForum, nil
	case sql.ErrNoRows:
		return nil, errors.ErrNotFoundInDB
	default:
		return nil, errors.ErrInternalError
	}
}
