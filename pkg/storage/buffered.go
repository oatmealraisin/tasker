package storage

import (
	"fmt"
	"os"

	"github.com/oatmealraisin/tasker/pkg/models"
)

type getZeroGuidError struct{}

func (e getZeroGuidError) Error() string {
	return "GUID 0 is reserved, cannot get"
}

type bufferStorage struct {
	buffer_guid   map[uint64]*models.Task
	buffer_tag    map[string][]uint64
	buffer_name   map[string][]uint64
	sort_priority []*models.Task
	queue         []models.Task
}

func newBufferStorage() *bufferStorage {
	return &bufferStorage{
		buffer_guid:   map[uint64]*models.Task{},
		buffer_tag:    map[string][]uint64{},
		buffer_name:   map[string][]uint64{},
		sort_priority: []*models.Task{},
		queue:         []models.Task{},
	}
}

func (b *bufferStorage) updateBuffers(p_t *models.Task) {
	b.buffer_guid[p_t.Guid] = p_t

	b.buffer_name[p_t.Name] = append(b.buffer_name[p_t.Name], p_t.Guid)

	for _, v := range p_t.Tags {
		b.buffer_tag[v] = append(b.buffer_tag[v], p_t.Guid)
	}
}

func (b *bufferStorage) DeleteTask(guid uint64) error {
	if _, ok := b.buffer_guid[guid]; !ok {
		return fmt.Errorf("bufferStorage.DeleteTask: Guid %d not found.", guid)
	}
	task, _ := b.buffer_guid[guid]

	delete(b.buffer_guid, guid)

	b.removeTaskFromTagBuffer(*task)
	b.removeTaskFromNameBuffer(*task)

	return nil
}

func (b *bufferStorage) removeTaskFromTagBuffer(task models.Task) {
	for _, tag := range task.Tags {
		if _, ok := b.buffer_tag[tag]; ok {
			b.buffer_tag[tag] = removeUuid(b.buffer_tag[tag], task.Guid)
		}
	}
}

func (b *bufferStorage) removeTaskFromNameBuffer(task models.Task) {
	if name_list, ok := b.buffer_name[task.Name]; ok {
		b.buffer_name[task.Name] = removeUuid(name_list, task.Guid)
	}
}

func removeUuid(l []uint64, u uint64) []uint64 {
	for i, uuid := range l {
		if u == uuid {
			l[len(l)-1], l[i] = l[i], l[len(l)-1]
			return l[:len(l)-1]
		}
	}

	return l
}

func (b *bufferStorage) GetAllTags() []string {
	result := make([]string, len(b.buffer_tag))

	i := 0
	for k := range b.buffer_tag {
		result[i] = k
		i++
	}

	return result
}

func (b *bufferStorage) GetByTags(tags []string) []uint64 {
	result := []uint64{}
	for _, tag := range tags {
		if tasks, ok := b.buffer_tag[tag]; ok {
			result = append(result, tasks...)
		} else {
			fmt.Fprintf(os.Stderr, "Could not find tasks with tag '%s'\n", tag)
		}
	}

	return result
}

func (b *bufferStorage) EditTask(oldTask, newTask models.Task) error {
	if oldTask.Guid != newTask.Guid {
		return fmt.Errorf("Cannot change the GUID of a Task.")
	}

	if oldTask.Added != newTask.Added {
		return fmt.Errorf("Cannot change the add date of a Task")
	}

	if oldTask.Name != newTask.Name {
		b.removeTaskFromNameBuffer(oldTask)
	}

	tagList := make(map[string]bool)
	for _, tag := range oldTask.Tags {
		tagList[tag] = false
	}

	for _, tag := range newTask.Tags {
		if _, ok := tagList[tag]; ok {
			delete(tagList, tag)
		} else {
			b.buffer_tag[tag] = append(b.buffer_tag[tag], newTask.Guid)
		}
	}

	for tag := range tagList {
		b.buffer_tag[tag] = removeUuid(b.buffer_tag[tag], oldTask.Guid)
	}

	b.buffer_guid[oldTask.Guid] = &newTask

	return nil
}
