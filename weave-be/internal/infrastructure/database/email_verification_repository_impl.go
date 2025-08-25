package database

import (
	"context"
	"time"

	"gorm.io/gorm"
	"weave-module/database"
	"weave-module/models"
	"weave-be/internal/domain/entities"
	"weave-be/internal/domain/repositories"
)

type emailVerificationRepositoryImpl struct {
	db *gorm.DB
}

func NewEmailVerificationRepository() repositories.EmailVerificationRepository {
	return &emailVerificationRepositoryImpl{
		db: database.GetDB(),
	}
}

// Convert between domain entity and database model
func (r *emailVerificationRepositoryImpl) entityToModel(verification *entities.EmailVerification) *models.EmailVerification {
	return &models.EmailVerification{
		ID:        verification.ID,
		Email:     verification.Email,
		Code:      verification.Code,
		ExpiresAt: verification.ExpiresAt,
		IsUsed:    verification.IsUsed,
		UserID:    verification.UserID,
		CreatedAt: verification.CreatedAt,
		UpdatedAt: verification.UpdatedAt,
	}
}

func (r *emailVerificationRepositoryImpl) modelToEntity(model *models.EmailVerification) *entities.EmailVerification {
	return &entities.EmailVerification{
		ID:        model.ID,
		Email:     model.Email,
		Code:      model.Code,
		ExpiresAt: model.ExpiresAt,
		IsUsed:    model.IsUsed,
		UserID:    model.UserID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

// Create operations
func (r *emailVerificationRepositoryImpl) Create(ctx context.Context, verification *entities.EmailVerification) error {
	model := r.entityToModel(verification)
	return r.db.WithContext(ctx).Create(model).Error
}

// Read operations
func (r *emailVerificationRepositoryImpl) GetByCode(ctx context.Context, code string) (*entities.EmailVerification, error) {
	var model models.EmailVerification
	err := r.db.WithContext(ctx).Where("code = ? AND is_used = false", code).First(&model).Error
	if err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *emailVerificationRepositoryImpl) GetActiveByEmail(ctx context.Context, email string) (*entities.EmailVerification, error) {
	var model models.EmailVerification
	err := r.db.WithContext(ctx).Where("email = ? AND is_used = false AND expires_at > ?", email, time.Now()).First(&model).Error
	if err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

// Update operations
func (r *emailVerificationRepositoryImpl) Update(ctx context.Context, verification *entities.EmailVerification) error {
	model := r.entityToModel(verification)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete operations
func (r *emailVerificationRepositoryImpl) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ? OR is_used = true", time.Now()).Delete(&models.EmailVerification{}).Error
}

func (r *emailVerificationRepositoryImpl) DeleteByEmail(ctx context.Context, email string) error {
	return r.db.WithContext(ctx).Where("email = ? AND is_used = false", email).Delete(&models.EmailVerification{}).Error
}

// Utility operations
func (r *emailVerificationRepositoryImpl) CountActiveByEmail(ctx context.Context, email string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.EmailVerification{}).
		Where("email = ? AND is_used = false AND expires_at > ?", email, time.Now()).
		Count(&count).Error
	return count, err
}