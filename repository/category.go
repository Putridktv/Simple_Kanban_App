package repository

import (
	"context"
	"errors"
	"kanbanApp/entity"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	var result []entity.Category
	err := r.db.WithContext(ctx).Table("categories").Where("user_id = ?", id).Find(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.Category{}, nil
		}
		return nil, err
	}

	return result, nil
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	err = r.db.WithContext(ctx).Create(&category).Error
	if err != nil {
		return 0, err
	}
	categoryId = category.ID
	return categoryId, nil
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	err := r.db.WithContext(ctx).Create(&categories)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	var result entity.Category
	err := r.db.WithContext(ctx).Table("categories").Where("id = ?", id).Find(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Category{}, nil
		}
		return entity.Category{}, err
	}
	return result, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	upCategory := r.db.WithContext(ctx).Table("categories").Where("id = ?", category.ID).Updates(&category)
	if upCategory != nil {
		return upCategory.Error
	}
	return nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	var result entity.Category
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&result)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
