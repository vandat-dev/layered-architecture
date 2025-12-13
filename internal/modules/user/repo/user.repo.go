package repo

import (
	"app/internal/modules/user/dto"
	"app/internal/modules/user/model"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(email string) *model.User
	GetUserByUsername(username string) *model.User
	GetUserByID(id uuid.UUID) *model.User
	GetListUser(req dto.UserListRequestDto) ([]*model.User, int64, error)
	CreateUser(user *model.User) (uuid.UUID, error)
	UpdateUser(id uuid.UUID, user *model.User) (*model.User, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) GetUserByEmail(email string) *model.User {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

func (r *userRepository) GetUserByUsername(username string) *model.User {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

func (r *userRepository) GetUserByID(id uuid.UUID) *model.User {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil
	}
	return &user
}

func (r *userRepository) GetListUser(req dto.UserListRequestDto) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.Model(&model.User{})

	if req.Email != "" {
		query = query.Where("email ILIKE ?", "%"+req.Email+"%")
	}

	if req.Username != "" {
		query = query.Where("username", req.Username)
	}

	if req.SystemRole != "" {
		query = query.Where("system_role", req.SystemRole)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Limit(req.Limit).Offset(req.Skip).Order("created_at DESC")

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) CreateUser(user *model.User) (uuid.UUID, error) {
	err := r.db.Create(user).Error
	return user.ID, err
}

func (r *userRepository) UpdateUser(id uuid.UUID, user *model.User) (*model.User, error) {
	// run update
	if err := r.db.Model(&model.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		return nil, err
	}

	var updatedUser model.User
	if err := r.db.First(&updatedUser, id).Error; err != nil {
		return nil, err
	}

	return &updatedUser, nil
}
