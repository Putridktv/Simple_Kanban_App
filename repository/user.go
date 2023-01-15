package repository

import (
	"context"
	"errors"
	"kanbanApp/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	var result entity.User
	err := r.db.WithContext(ctx).Table("users").Where("id = ?", id).Find(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, nil
		}
		return entity.User{}, err
	}
	return result, nil // TODO: replace this
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var result entity.User
	err := r.db.WithContext(ctx).Table("users").Where("email = ?", email).Find(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, nil
		}
		return entity.User{}, err
	}
	return result, nil // TODO: replace this
}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.WithContext(ctx).Create(&user)
	if err.Error != nil {
		return entity.User{}, err.Error
	}
	return user, nil // TODO: replace this
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	upUser := r.db.WithContext(ctx).Table("users").Where("id = ?", user.ID).Updates(user)
	if upUser != nil {
		return entity.User{}, upUser.Error
	}
	return user, nil // TODO: replace this
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	var result entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&result)
	if err.Error != nil {
		return err.Error
	}
	return nil // TODO: replace this
}
