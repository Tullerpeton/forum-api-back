package admin

import "github.com/forum-api-back/internal/pkg/models"

type UseCase interface {
	ClearBase() error
	GetBaseDetails() (*models.BaseDetails, error)
}
