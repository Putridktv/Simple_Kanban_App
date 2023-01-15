package repository

import (
	"context"
	"errors"
	"kanbanApp/entity"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	var result []entity.Task
	err := r.db.WithContext(ctx).Table("tasks").Where("user_id = ?", id).Find(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.Task{}, nil
		}
		return nil, err
	}
	return result, nil
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	err = r.db.WithContext(ctx).Create(&task).Error
	if err != nil {
		return 0, err
	}
	taskId = task.ID
	return taskId, nil
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	var result entity.Task
	err := r.db.WithContext(ctx).Table("tasks").Where("id = ?", id).Find(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Task{}, nil
		}
		return entity.Task{}, err
	}
	return result, nil
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	var result []entity.Task
	err := r.db.WithContext(ctx).Table("tasks").Where("category_id = ?", catId).Find(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.Task{}, nil
		}
		return nil, err
	}
	return result, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	upTask := r.db.WithContext(ctx).Table("tasks").Where("id = ?", task.ID).Updates(&task)
	if upTask != nil {
		return upTask.Error
	}
	return nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	var result entity.Task
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&result)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
