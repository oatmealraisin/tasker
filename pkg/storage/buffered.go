package storage

import (
	"fmt"

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

func (b *bufferStorage) DeleteTask(guid uint64) error {
	if _, ok := b.buffer_guid[guid]; !ok {
		return fmt.Errorf("bufferStorage.DeleteTask: Guid %d not found.", guid)
	}
	task, _ := b.buffer_guid[guid]

	delete(b.buffer_guid, guid)

	for _, tag := range task.Tags {
		if _, ok := b.buffer_tag[tag]; !ok {
			fmt.Println("Hmmm 1")
			continue
		}

		tag_list, _ := b.buffer_tag[tag]
		loc := -1
		for i, tag_task := range tag_list {
			if task.Guid == tag_task.Guid {
				loc = i
				break
			}
		}
		if loc == -1 {
			fmt.Println("Hmmmm 2")
		} else {
			tag_list[len(tag_list)-1], tag_list[loc] = tag_list[loc], tag_list[len(tag_list)-1]
			b.buffer_tag[tag] = tag_list[:len(tag_list)-1]
		}
	}

	if _, ok := b.buffer_name[task.Name]; !ok {
		return fmt.Errorf("Hmmm3")
	}

	name_list, _ := b.buffer_name[task.Name]
	loc := -1
	for i, name_task := range name_list {
		if task.Guid == name_task.Guid {
			loc = i
			break
		}
	}
	if loc == -1 {
		return fmt.Errorf("Hmmm4")
	} else {
		name_list[len(name_list)-1], name_list[loc] = name_list[loc], name_list[len(name_list)-1]
		b.buffer_name[task.Name] = name_list[:len(name_list)-1]
	}

	return nil
}
