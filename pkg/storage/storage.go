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
	GetByTag(tag string) []uint64
	GetByName(name string) []uint64
	GetAllTasks() []uint64
	GetAllTags() []string

	DeleteTask(guid uint64) error
}

func setupStorageDir() error {
	wd := viper.GetString("WorkingDir")

	err := os.MkdirAll(wd, 0744)
	return err
}

type GetFunc func(uuid uint64) (models.Task, error)
