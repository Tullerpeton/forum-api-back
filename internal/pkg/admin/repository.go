package admin

import "github.com/forum-api-back/internal/pkg/models"

type Repository interface {
	ClearBase() error
	SelectBaseDetails() (*models.BaseDetails, error)
}
