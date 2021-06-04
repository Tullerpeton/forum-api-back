package repository

import (
	"database/sql"

	"github.com/forum-api-back/internal/pkg/models"
	"github.com/forum-api-back/internal/pkg/user"
	"github.com/forum-api-back/pkg/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) user.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

func (r *PostgresqlRepository) InsertUser(userInfo *models.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users(nickname, fullname, about, email) "+
			"VALUES ($1, $2, $3, $4)",
		userInfo.NickName,
		userInfo.FullName,
		userInfo.About,
		userInfo.Email,
	)

	if err != nil {
		return errors.ErrDataConflict
	}

	return nil
}

func (r *PostgresqlRepository) SelectUserByEmailOrNickname(email, nickname string) (*models.User, error) {
	row := r.db.QueryRow(
		"SELECT nickname, fullname, about, email "+
			"FROM users "+
			"WHERE nickname = $1 OR email = $2",
		nickname,
		email,
	)

	selectedUser := &models.User{}
	about := sql.NullString{}
	err := row.Scan(
		&selectedUser.NickName,
		&selectedUser.FullName,
		&about,
		&selectedUser.Email,
	)
	selectedUser.About = about.String

	switch err {
	case nil:
		return selectedUser, nil
	case sql.ErrNoRows:
		return nil, errors.ErrNotFoundInDB
	default:
		return nil, errors.ErrInternalError
	}
}

func (r *PostgresqlRepository) SelectUserByNickName(nickname string) (*models.User, error) {
	row := r.db.QueryRow(
		"SELECT nickname, fullname, about, email "+
			"FROM users "+
			"WHERE nickname = $1",
		nickname,
	)

	selectedUser := &models.User{}
	about := sql.NullString{}
	err := row.Scan(
		&selectedUser.NickName,
		&selectedUser.FullName,
		&about,
		&selectedUser.Email,
	)
	selectedUser.About = about.String

	switch err {
	case nil:
		return selectedUser, nil
	case sql.ErrNoRows:
		return nil, errors.ErrNotFoundInDB
	default:
		return nil, errors.ErrInternalError
	}
}

func (r *PostgresqlRepository) SelectUsersByForum(forumSlug string) ([]*models.User, error) {
	rows, err := r.db.Query(
		"SELECT nickname, fullname, about, email "+
			"FROM users u "+
			"JOIN threads t ON (u.nickname = t.author_nickname AND t.forum_slug = $1) "+
			"JOIN posts p ON (u.nickname = p.author_nickname AND p.forum_slug = $1) "+
			"GROUP BY nickname "+
			"ORDER BY nickname",
		forumSlug,
	)

	switch err {
	case nil:
		users := make([]*models.User, 0)
		for rows.Next() {
			selectedUser := &models.User{}
			err := rows.Scan(
				&selectedUser.NickName,
				&selectedUser.FullName,
				&selectedUser.About,
				&selectedUser.Email,
			)
			if err != nil {
				return nil, errors.ErrInternalError
			}

			users = append(users, selectedUser)
		}
		return users, nil
	case sql.ErrNoRows:
		return nil, errors.ErrNotFoundInDB
	default:
		return nil, errors.ErrInternalError
	}
}

func (r *PostgresqlRepository) UpdateUserProfile(userInfo *models.User) error {
	row := r.db.QueryRow(
		"SELECT email "+
			"FROM users "+
			"WHERE nickname = $1",
		userInfo.NickName,
	)
	var email string
	err := row.Scan(&email)
	if err != nil {
		return errors.ErrNotFoundInDB
	}
	if email == userInfo.Email {
		return errors.ErrDataConflict
	}

	_, err = r.db.Exec(
		"UPDATE users SET "+
			"email = $1, "+
			"fullname = $2, "+
			"about = $3 "+
			"WHERE nickname = $4",
		userInfo.Email,
		userInfo.FullName,
		userInfo.About,
		userInfo.NickName,
	)

	if err != nil {
		return errors.ErrDataConflict
	}

	return nil
}
