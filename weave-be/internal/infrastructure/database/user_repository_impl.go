package database

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"weave-module/database"
	"weave-module/models"
	"weave-be/internal/domain/entities"
	"weave-be/internal/domain/repositories"
)

// userRepositoryImpl implements the UserRepository interface
// This follows the Repository Pattern and Adapter Pattern
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository implementation
func NewUserRepository() repositories.UserRepository {
	return &userRepositoryImpl{
		db: database.GetDB(),
	}
}

// Convert between domain entity and database model
func (r *userRepositoryImpl) entityToModel(user *entities.User) *models.User {
	return &models.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		ProfileImage: user.ProfileImage,
		Bio:          user.Bio,
		IsVerified:   user.IsVerified,
		IsActive:     user.IsActive,
		GoogleID:     user.GoogleID,
		GoogleEmail:  user.GoogleEmail,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func (r *userRepositoryImpl) modelToEntity(model *models.User) *entities.User {
	return &entities.User{
		ID:           model.ID,
		Username:     model.Username,
		Email:        model.Email,
		PasswordHash: model.PasswordHash,
		ProfileImage: model.ProfileImage,
		Bio:          model.Bio,
		IsVerified:   model.IsVerified,
		IsActive:     model.IsActive,
		GoogleID:     model.GoogleID,
		GoogleEmail:  model.GoogleEmail,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}

func (r *userRepositoryImpl) modelsToEntities(models []*models.User) []*entities.User {
	entities := make([]*entities.User, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(model)
	}
	return entities
}

// Create operations
func (r *userRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	model := r.entityToModel(user)
	return r.db.WithContext(ctx).Create(model).Error
}

// Read operations
func (r *userRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var model models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var model models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *userRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	var model models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&model).Error
	if err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *userRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	var models []*models.User
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *userRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}

// Update operations
func (r *userRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	model := r.entityToModel(user)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *userRepositoryImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, profileImage *string, bio *string) error {
	updates := make(map[string]interface{})
	if profileImage != nil {
		updates["profile_image"] = *profileImage
	}
	if bio != nil {
		updates["bio"] = *bio
	}
	
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *userRepositoryImpl) UpdateVerificationStatus(ctx context.Context, userID uuid.UUID, isVerified bool) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("is_verified", isVerified).Error
}

func (r *userRepositoryImpl) UpdateActiveStatus(ctx context.Context, userID uuid.UUID, isActive bool) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("is_active", isActive).Error
}

// Delete operations
func (r *userRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *userRepositoryImpl) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("is_active", false).Error
}

// Search operations
func (r *userRepositoryImpl) SearchByUsername(ctx context.Context, query string, limit, offset int) ([]*entities.User, error) {
	var models []*models.User
	err := r.db.WithContext(ctx).
		Where("username ILIKE ? AND is_active = ?", "%"+query+"%", true).
		Order("username ASC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *userRepositoryImpl) SearchByUsernameCount(ctx context.Context, query string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("username ILIKE ? AND is_active = ?", "%"+query+"%", true).
		Count(&count).Error
	return count, err
}

func (r *userRepositoryImpl) SearchByEmail(ctx context.Context, query string, limit, offset int) ([]*entities.User, error) {
	var models []*models.User
	err := r.db.WithContext(ctx).
		Where("email ILIKE ? AND is_active = ?", "%"+query+"%", true).
		Order("email ASC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

// Existence checks
func (r *userRepositoryImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *userRepositoryImpl) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// Follow system
func (r *userRepositoryImpl) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	follow := &models.UserFollow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	return r.db.WithContext(ctx).Create(follow).Error
}

func (r *userRepositoryImpl) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&models.UserFollow{}).Error
}

func (r *userRepositoryImpl) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.UserFollow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

func (r *userRepositoryImpl) GetFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.User, error) {
	var models []*models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN user_follows ON users.id = user_follows.follower_id").
		Where("user_follows.following_id = ? AND users.is_active = ?", userID, true).
		Order("user_follows.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *userRepositoryImpl) GetFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.User, error) {
	var models []*models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN user_follows ON users.id = user_follows.following_id").
		Where("user_follows.follower_id = ? AND users.is_active = ?", userID, true).
		Order("user_follows.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *userRepositoryImpl) GetFollowersCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.UserFollow{}).
		Where("following_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *userRepositoryImpl) GetFollowingCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.UserFollow{}).
		Where("follower_id = ?", userID).
		Count(&count).Error
	return count, err
}

// OAuth operations
func (r *userRepositoryImpl) GetByGoogleID(ctx context.Context, googleID string) (*entities.User, error) {
	var model models.User
	err := r.db.WithContext(ctx).Where("google_id = ?", googleID).First(&model).Error
	if err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *userRepositoryImpl) ExistsByGoogleID(ctx context.Context, googleID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("google_id = ?", googleID).
		Count(&count).Error
	return count > 0, err
}