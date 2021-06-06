package repository

import (
	"database/sql"
	"fmt"
	"strings"

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

func (r *PostgresqlRepository) SelectUserByEmailOrNickname(email, nickname string) ([]*models.User, error) {
	rows, err := r.db.Query(
		"SELECT nickname, fullname, about, email "+
			"FROM users "+
			"WHERE nickname = $1 OR email = $2",
		nickname,
		email,
	)

	if err != nil {
		return nil, errors.ErrNotFoundInDB
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		selectedUser := &models.User{}
		about := sql.NullString{}
		err := rows.Scan(
			&selectedUser.NickName,
			&selectedUser.FullName,
			&about,
			&selectedUser.Email,
		)
		selectedUser.About = about.String

		if err != nil {
			return nil, errors.ErrInternalError
		}

		users = append(users, selectedUser)
	}

	return users, nil
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

func (r *PostgresqlRepository) SelectUsersByForum(forumSlug string,
	paginator *models.UserPaginator) ([]*models.User, error) {
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
	if paginator.Since == "" {
		rows, err = r.db.Query(
			"SELECT u.nickname, u.fullname, u.about, u.email "+
				"FROM users u "+
				"JOIN authors a ON (u.nickname = a.user_nickname AND a.forum_slug = $1) "+
				"ORDER BY u.nickname "+orderSort+
				"LIMIT $2",
			forumSlug,
			paginator.Limit,
		)
	} else {
		rows, err = r.db.Query(
			"SELECT u.nickname, u.fullname, u.about, u.email "+
				"FROM users u "+
				"JOIN authors a ON (u.nickname = a.user_nickname AND a.forum_slug = $1) "+
				"WHERE (u.nickname"+orderCompare+"$2) "+
				"ORDER BY u.nickname "+orderSort+
				"LIMIT $3",
			forumSlug,
			paginator.Since,
			paginator.Limit,
		)
	}

	switch err {
	case nil:
		defer rows.Close()
		users := make([]*models.User, 0)
		for rows.Next() {
			about := sql.NullString{}
			selectedUser := &models.User{}
			err := rows.Scan(
				&selectedUser.NickName,
				&selectedUser.FullName,
				&about,
				&selectedUser.Email,
			)
			selectedUser.About = about.String
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
	columns := make([]string, 0)
	args := make([]interface{}, 1)
	args[0] = userInfo.NickName
	if userInfo.Email != "" {
		args = append(args, userInfo.Email)
		columns = append(columns, fmt.Sprintf("email = $%d ", len(args)))
	}
	if userInfo.FullName != "" {
		args = append(args, userInfo.FullName)
		columns = append(columns, fmt.Sprintf("fullname = $%d ", len(args)))
	}
	if userInfo.About != "" {
		args = append(args, userInfo.About)
		columns = append(columns, fmt.Sprintf("about = $%d ", len(args)))
	}

	if len(columns) == 0 {
		return nil
	}

	_, err := r.db.Exec(
		"UPDATE users SET "+
			strings.Join(columns, ", ")+
			" WHERE nickname = $1",
		args...,
	)

	if err != nil {
		return errors.ErrDataConflict
	}

	return nil
}
