package models

import (
	"fmt"
	"os"
	"sort"

	"github.com/oatmealraisin/tasker/pkg/gui/terminal"
)

// TaskTree is a data structure for navigating the parent/child relationship of
// tasks.
type TaskTree struct {
	nodes map[uint64]*TaskTree
	key   uint64
}

// CreateTaskTree takes a list of tasks and fully expands every tree touched by
// that list.
func CreateTaskTree(tasks []uint64, get func(uuid uint64) (Task, error)) *TaskTree {
	result := &TaskTree{}
	result.nodes = make(map[uint64]*TaskTree)

	for _, uuid := range tasks {
		result.Insert(uuid, get)
	}

	return result
}

// Insert finds where the given task would fit into the list of trees. If it
// doesn't fit anywhere, it gets added to the top level (along with all parents
// and children)
func (tt *TaskTree) Insert(uuid uint64, get func(uuid uint64) (Task, error)) error {
	_, err := tt.recurseInsert(uuid, get)

	return err
}

func (tt *TaskTree) recurseInsert(uuid uint64, get func(uuid uint64) (Task, error)) (*TaskTree, error) {
	task, err := get(uuid)
	if err != nil {
		return nil, err
	}

	if task.Parent == tt.key {
		if _, ok := tt.nodes[uuid]; !ok {
			tt.nodes[uuid] = &TaskTree{
				nodes: make(map[uint64]*TaskTree),
				key:   uuid,
			}

			return tt.nodes[uuid], nil
		}

		return tt.nodes[uuid], nil
	}

	sub, err := tt.recurseInsert(task.Parent, get)
	if err != nil {
		return nil, err
	}

	if sub.key == task.Parent {
		if _, ok := sub.nodes[uuid]; !ok {
			sub.nodes[uuid] = &TaskTree{
				nodes: make(map[uint64]*TaskTree),
				key:   uuid,
			}

			return sub.nodes[uuid], nil
		}

		return sub.nodes[uuid], nil
	}

	return nil, fmt.Errorf("Couldn't insert node %d\n", uuid)
}

// Print displays the contents of the tree up to the third level. Anything
// further will be displayed with ellipses.
func (tt *TaskTree) Print(get func(uuid uint64) (Task, error)) {
	w := terminal.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)

	printColumns(w)

	roots := make([]uint64, len(tt.nodes))
	i := 0
	for k := range tt.nodes {
		roots[i] = k
		i++
	}
	sort.Slice(roots, func(i, j int) bool { return roots[i] < roots[j] })

	for _, key := range roots {
		root := tt.nodes[key]
		root._print([]rune{}, w, get)
	}

	w.Flush()
}

func (tt *TaskTree) _print(prefix []rune, w *terminal.Writer, get func(uuid uint64) (Task, error)) {
	// Start with printing the root task
	task, err := get(tt.key)
	if err != nil {
		// TODO: Log
		fmt.Println("Could not get task %d: %s\n", tt.key, err.Error())
		return
	}

	sfyTask := task.stringify()

	// If we're not a true root, we need to add some lines
	if len(prefix) > 0 {
		if len(tt.nodes) > 0 {
			sfyTask["name"] = fmt.Sprintf("%s%s> %s", string(prefix), "┬─", sfyTask["name"])
		} else {
			sfyTask["name"] = fmt.Sprintf("%s%s> %s", string(prefix), "──", sfyTask["name"])
		}
	}

	printStringyTask(w, sfyTask)

	subs := make([]uint64, len(tt.nodes))
	i := 0
	for k := range tt.nodes {
		subs[i] = k
		i++
	}
	sort.Slice(subs, func(i, j int) bool { return subs[i] < subs[j] })

	final := len(tt.nodes) - 1

	prefix = append(prefix, '├')
	for i, skey := range subs {
		// This needs to be inside the loop, so that we can revert the inner
		// recursion's changes
		if i == final {
			prefix[len(prefix)-1] = '└'
		} else {
			prefix[len(prefix)-1] = '├'
		}

		// If we're grandchildren, we need to change the previous prefix
		if len(prefix) > 1 {
			switch prefix[len(prefix)-2] {
			case '└':
				prefix[len(prefix)-2] = ' '
				break
			case '├':
				prefix[len(prefix)-2] = '│'
				break
			}
		}

		sub := tt.nodes[skey]
		sub._print(prefix, w, get)
	}
}
