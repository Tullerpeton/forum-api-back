package repository

import (
	"database/sql"

	"github.com/forum-api-back/internal/pkg/admin"
	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/pkg/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) admin.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

func (r *PostgresqlRepository) ClearBase() error {
	_, err := r.db.Exec("TRUNCATE  votes, posts, threads, forums, users CASCADE")

	if err != nil {
		return errors.ErrInternalError
	}
	return nil
}

func (r *PostgresqlRepository) SelectBaseDetails() (*models.BaseDetails, error) {
	baseDetails := &models.BaseDetails{}
	row := r.db.QueryRow(
		"SELECT " +
			"(SELECT COUNT(*) FROM forums) AS forums, " +
			"(SELECT COUNT(*) FROM threads) AS threads, " +
			"(SELECT COUNT(*) FROM posts) AS posts, " +
			"(SELECT COUNT(*) FROM users) AS users",
	)

	err := row.Scan(
		&baseDetails.Forum,
		&baseDetails.Thread,
		&baseDetails.Post,
		&baseDetails.User,
	)

	if err != nil {
		return nil, errors.ErrInternalError
	}

	return baseDetails, err
}
