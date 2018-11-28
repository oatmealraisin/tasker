package storage

import (
	"os"

	"github.com/oatmealraisin/tasker/pkg/models"
)

type SqliteStorage struct {
	f      os.File
	buffer map[uint64]*models.Task
}

func (s *SqliteStorage) CreateTask(t *models.Task) error {
	panic("not implemented")
}

func (s *SqliteStorage) CreateTasks(t []*models.Task) []error {
	panic("not implemented")
}

func (s *SqliteStorage) GetTask(guid uint64) (models.Task, error) {
	panic("not implemented")
}

func (s *SqliteStorage) GetByTag(tag string) ([]models.Task, error) {
	panic("not implemented")
}

func (s *SqliteStorage) GetByName(name string) ([]models.Task, error) {
	panic("not implemented")
}

func (s *SqliteStorage) GetAllTasks() ([]models.Task, error) {
	panic("not implemented")
}
