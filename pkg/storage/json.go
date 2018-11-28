package storage

import (
	"os"

	"github.com/oatmealraisin/tasker/pkg/models"
)

type JsonStorage struct {
	f os.File
}

func (j *JsonStorage) CreateTask(t *models.Task) error {
	panic("not implemented")
}

func (j *JsonStorage) CreateTasks(t []*models.Task) []error {
	panic("not implemented")
}

func (j *JsonStorage) GetTask(guid uint64) (models.Task, error) {
	panic("not implemented")
}

func (j *JsonStorage) GetByTag(tag string) ([]models.Task, error) {
	panic("not implemented")
}

func (j *JsonStorage) GetByName(name string) ([]models.Task, error) {
	panic("not implemented")
}

func (j *JsonStorage) GetAllTasks() ([]models.Task, error) {
	panic("not implemented")
}
