package storage

import (
	"os"

	"github.com/oatmealraisin/tasker/pkg/models"
	"github.com/spf13/viper"
)

type Storage interface {
	CreateTask(t models.Task) error
	CreateTasks(t []models.Task) []error
	GetTask(guid uint64) (models.Task, error)
	GetByTag(tag string) ([]models.Task, error)
	GetByName(name string) ([]models.Task, error)
	GetAllTasks() ([]models.Task, error)
}

func setupStorageDir() error {
	wd := viper.GetString("WorkingDir")

	err := os.MkdirAll(wd, 0744)
	return err
}
