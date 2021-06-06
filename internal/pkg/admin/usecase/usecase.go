package usecase

import (
	"github.com/forum-api-back/internal/pkg/admin"
	"github.com/forum-api-back/internal/pkg/models"
)

type AdminUseCase struct {
	AdminRepo admin.Repository
}

func NewUseCase(adminRepo admin.Repository) admin.UseCase {
	return &AdminUseCase{
		AdminRepo: adminRepo,
	}
}

func (u *AdminUseCase) ClearBase() error {
	return u.AdminRepo.ClearBase()
}

func (u *AdminUseCase) GetBaseDetails() (*models.BaseDetails, error) {
	return u.AdminRepo.SelectBaseDetails()
}
