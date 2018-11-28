package storage

import "github.com/oatmealraisin/tasker/pkg/models"

type PostgresStorage struct{}

func (p *PostgresStorage) CreateTask(t *models.Task) error {
	panic("not implemented")
}

func (p *PostgresStorage) CreateTasks(t []*models.Task) []error {
	panic("not implemented")
}

func (p *PostgresStorage) GetTask(guid uint64) (models.Task, error) {
	panic("not implemented")
}

func (p *PostgresStorage) GetByTag(tag string) ([]models.Task, error) {
	panic("not implemented")
}

func (p *PostgresStorage) GetByName(name string) ([]models.Task, error) {
	panic("not implemented")
}

func (p *PostgresStorage) GetAllTasks() ([]models.Task, error) {
	panic("not implemented")
}
