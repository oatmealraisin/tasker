package storage

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/oatmealraisin/tasker/pkg/models"
)

type CsvStorage struct {
	*bufferStorage
	f *os.File
}

func NewCsvStorage(filename string) Storage {
	var err error

	err = setupStorageDir()
	if err != nil {
		return nil
	}

	newDb, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil
	}

	result := new(CsvStorage)
	result.bufferStorage = newBufferStorage()
	result.f = newDb

	result.loadTasks(10000)

	return result
}

func (c *CsvStorage) GetTask(guid uint64) (models.Task, error) {
	if guid == 0 {
		return models.Task{}, getZeroGuidError{}
	}

	if _, ok := c.buffer_guid[guid]; !ok {
		return c.getTaskFromFile(guid)
	}

	return *c.buffer_guid[guid], nil
}

func (c *CsvStorage) GetByTag(tag string) []uint64 {
	panic("not implemented")
}

func (s *CsvStorage) GetByName(name string) []uint64 {
	result := make([]uint64, len(s.buffer_name[name]))
	if len(result) == 0 {
		return nil
	}

	for i, task := range s.buffer_name[name] {
		result[i] = task.Guid
	}

	return result
}

func (c *CsvStorage) GetAllTasks() []uint64 {
	// NOTE: For now, we load everything into memory
	var result []uint64
	var uuid uint64

	_, err := c.f.Seek(0, 0)
	if err != nil {
		fmt.Printf("Error resetting task file: %s\n", err.Error())
	}

	r := csv.NewReader(c.f)

	for {
		record, err := r.Read()
		if err == io.EOF {
			return result
		}

		uuid, err = strconv.ParseUint(record[0], 10, 64)
		if err != nil {
			// TODO: Log
			fmt.Printf("Invalid UUID: %s", record[0])
		}

		result = append(result, uuid)
	}

	return result
}

func (c *CsvStorage) CreateTask(t models.Task) error {
	if task, err := c.GetTask(t.Guid); err == nil {
		return fmt.Errorf("Task with GUID %d already exists:\n\t%s\n", t.Guid, task.Name)
	} else {
		if _, ok := err.(getZeroGuidError); !ok {
			return err
		}

		t.Guid = c.getNextGuid()
	}

	if t.Added == nil {
		t.Added = ptypes.TimestampNow()
	}

	c.queue = append(c.queue, t)

	p_t := &c.queue[len(c.queue)-1]

	c.updateBuffers(p_t)

	if err := c.writeAll(); err != nil {
		return fmt.Errorf("CsvStorage.CreateTask: %s", err.Error())
	}

	return nil
}

func (c *CsvStorage) CreateTasks(t []models.Task) []error {
	panic("not implemented")
}

// TODO: This is terrible
func (c *CsvStorage) getNextGuid() uint64 {
	keys := make([]uint64, 0, len(c.buffer_guid))
	for k := range c.buffer_guid {
		keys = append(keys, k)
	}

	models.UuidSort(keys)

	return keys[len(keys)-1] + 1
}

func (s *CsvStorage) loadTasks(num int) {
	s.queue = []models.Task{}

	r := csv.NewReader(s.f)

	for i := 0; i < num; i++ {
		record, err := r.Read()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading CSV Storage: %s", err.Error())
			}

			return
		}

		newTask, err := TaskFromCsv(record)
		if err != nil {
			fmt.Printf("Error extracting task from CSV record: %s", err.Error())
			continue
		}

		s.queue = append(s.queue, newTask)
		p_t := &s.queue[len(s.queue)-1]
		s.updateBuffers(p_t)
	}
}

// TODO: This is terrible
func (c *CsvStorage) writeAll() error {

	if err := c.f.Truncate(0); err != nil {
		return fmt.Errorf("CsvStorage.writeAll: %s", err)
	}
	if _, err := c.f.Seek(0, 0); err != nil {
		return fmt.Errorf("CsvStorage.writeAll: %s", err)
	}

	keys := []uint64{}
	for k, _ := range c.buffer_guid {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, k := range keys {
		if _, err := c.f.WriteString(TaskToCSV(*c.buffer_guid[k])); err != nil {
			return fmt.Errorf("CsvStorage.writeAll: %s", err)
		}
	}

	return nil
}

func (s *CsvStorage) getTaskFromFile(guid uint64) (models.Task, error) {
	// NOTE: For now, we load everything into memory, so if we don't have a guid
	// in the buffer, it doesn't exist
	return models.Task{}, fmt.Errorf("CSV Storage Error: guid not found %d\n", guid)
}

func TaskToCSV(task models.Task) string {
	var subtasks string
	for _, guid := range task.Subtasks {
		if subtasks == "" {
			subtasks = strconv.FormatUint(guid, 10)
		} else {
			subtasks = strings.Join([]string{subtasks, strconv.FormatUint(guid, 10)}, "|")
		}
	}

	var depends string
	for _, guid := range task.Dependencies {
		if depends == "" {
			depends = strconv.FormatUint(guid, 10)
		} else {
			depends = strings.Join([]string{depends, strconv.FormatUint(guid, 10)}, "|")
		}
	}

	var added string
	if time_added, err := ptypes.Timestamp(task.Added); err == nil {
		added = time_added.Format("2006-01-02")
	}

	var finished string
	if time_finished, err := ptypes.Timestamp(task.Finished); err == nil {
		finished = time_finished.Format("2006-01-02")
	}

	var due string
	if time_due, err := ptypes.Timestamp(task.Due); err == nil {
		due = time_due.Format("2006-01-02")
	}

	result := fmt.Sprintf("%d,%s,%d,%s,%s,%s,%t,%t,%s,%d,%s,%d,%s,%s\n",
		task.Guid,
		task.Name,
		task.Size,
		added,
		finished,
		due,
		task.Removed,
		task.Repeats,
		strings.Join(task.Tags, "|"),
		task.Priority,
		task.Url,
		task.Parent,
		subtasks,
		depends,
	)

	return result
}

func TaskFromCsv(record []string) (models.Task, error) {
	newTask := models.Task{}
	var err error

	if len(record) != 14 {
		return newTask, fmt.Errorf("CSV line doesn't have right number of columns.\n%s\n", record)
	}

	guid, err := strconv.Atoi(record[0])
	if err != nil {
		return newTask, fmt.Errorf("TaskFromCSV: Could not extract guid: %s\n", err)
	}

	var size int
	if record[2] != "" {
		size, err = strconv.Atoi(record[2])
		if err != nil {
			return newTask, fmt.Errorf("TaskFromCSV: Could not extract size: %s\n", err)
		}
	}

	tAdded, err := time.Parse("2006-01-02", record[3])
	if err != nil {
		return newTask, err
	}

	added, err := ptypes.TimestampProto(tAdded)
	if err != nil {
		return newTask, err
	}

	var finished *timestamp.Timestamp
	if record[4] != "" {
		tFinished, err := time.Parse("2006-01-02", record[4])
		if err != nil {
			return newTask, err
		}

		finished, err = ptypes.TimestampProto(tFinished)
		if err != nil {
			return newTask, err
		}
	} else {
		finished = nil
	}

	var due *timestamp.Timestamp
	if record[5] != "" {
		tDue, err := time.Parse("2006-01-02", record[5])
		if err != nil {
			return newTask, err
		}

		due, err = ptypes.TimestampProto(tDue)
		if err != nil {
			return newTask, err
		}
	} else {
		due = nil
	}

	removed, err := strconv.ParseBool(record[6])
	if err != nil {
		return newTask, err
	}

	repeats, err := strconv.ParseBool(record[7])
	if err != nil {
		return newTask, err
	}

	var priority int
	if record[9] != "" {
		priority, err = strconv.Atoi(record[9])
		if err != nil {
			return newTask, fmt.Errorf("TaskFromCSV: Could not extract priority: %s\n", err)
		}
	}

	var parent uint64
	if record[11] != "" {
		if p, err := strconv.ParseUint(record[11], 10, 64); err == nil {
			parent = uint64(p)
		}
	}

	subtasks := []uint64{}
	if record[12] != "" {
		strSubTasks := strings.Split(record[12], "|")
		for _, strSubTask := range strSubTasks {
			subtask, err := strconv.Atoi(strSubTask)
			if err != nil {
				return newTask, fmt.Errorf("TaskFromCSV: Could not extract subtasks: %s\n", err)
			}

			subtasks = append(subtasks, uint64(subtask))
		}
	}

	depends := []uint64{}
	if record[13] != "" {
		strDepends := strings.Split(record[13], "|")
		for _, strDepend := range strDepends {
			depend, err := strconv.Atoi(strDepend)
			if err != nil {
				return newTask, fmt.Errorf("TaskFromCSV: Could not extract dependencies: %s\n", err)
			}

			depends = append(depends, uint64(depend))
		}
	}

	newTask = models.Task{
		Guid:         uint64(guid),
		Name:         record[1],
		Size:         uint32(size),
		Added:        added,
		Finished:     finished,
		Due:          due,
		Removed:      removed,
		Repeats:      repeats,
		Tags:         strings.Split(record[8], "|"),
		Priority:     uint32(priority),
		Url:          record[10],
		Parent:       parent,
		Subtasks:     subtasks,
		Dependencies: depends,
	}

	return newTask, nil
}
