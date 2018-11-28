package storage

import (
	"os"

	"github.com/oatmealraisin/tasker/pkg/models"
)

type YamlStorage struct {
	f os.File
}

func (y *YamlStorage) CreateTask(t *models.Task) error {
	panic("not implemented")
}

func (y *YamlStorage) CreateTasks(t []*models.Task) []error {
	panic("not implemented")
}

func (y *YamlStorage) GetTask(guid uint64) (models.Task, error) {
	panic("not implemented")
}

func (y *YamlStorage) GetByTag(tag string) ([]models.Task, error) {
	panic("not implemented")
}

func (y *YamlStorage) GetByName(name string) ([]models.Task, error) {
	panic("not implemented")
}

func (y *YamlStorage) GetAllTasks() ([]models.Task, error) {
	panic("not implemented")
}
