package storage

import (
	"github.com/oatmealraisin/tasker/pkg/models"
)

type bufferStorage struct {
	buffer_guid   map[uint64]*models.Task
	buffer_tag    map[string][]*models.Task
	buffer_name   map[string][]*models.Task
	sort_priority []*models.Task
	queue         []models.Task
}

func newBufferStorage() *bufferStorage {
	return &bufferStorage{
		buffer_guid:   map[uint64]*models.Task{},
		buffer_tag:    map[string][]*models.Task{},
		buffer_name:   map[string][]*models.Task{},
		sort_priority: []*models.Task{},
		queue:         []models.Task{},
	}
}

func (b *bufferStorage) updateBuffers(p_t *models.Task) {
	b.buffer_guid[p_t.Guid] = p_t

	b.buffer_name[p_t.Name] = append(b.buffer_name[p_t.Name], p_t)

	for _, v := range p_t.Tags {
		b.buffer_tag[v] = append(b.buffer_tag[v], p_t)
	}
}
