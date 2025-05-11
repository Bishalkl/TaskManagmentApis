package repositories

import "gorm.io/gorm"

type TaskRepository interface {
}

type TaskRepositoryImpl struct {
	DB *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &TaskRepositoryImpl{
		DB: db,
	}
}

// create
// update
// retrieve
// delete
