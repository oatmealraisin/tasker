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
		task, err := get(key)
		if err != nil {
			// TODO: Log
			fmt.Println("Could not get task %d: %s\n", root.key, err.Error())
			return
		}

		printStringyTask(w, task.stringify())

		subs := make([]uint64, len(root.nodes))
		i := 0
		for k := range root.nodes {
			subs[i] = k
			i++
		}
		sort.Slice(subs, func(i, j int) bool { return subs[i] < subs[j] })

		i = 0
		for _, skey := range subs {
			sub := root.nodes[skey]
			subTask, err := get(skey)
			// TODO: Log
			if err != nil {
				fmt.Println("Could not get task %d: %s\n", sub.key, err.Error())
			}

			sfyTask := subTask.stringify()
			if len(sub.nodes) > 0 {
				prefix := " "
				if len(root.nodes) > 1 && len(root.nodes) != i+1 {
					sfyTask["name"] = fmt.Sprintf("├┬─> %s", sfyTask["name"])
					prefix = "│"
				} else {
					sfyTask["name"] = fmt.Sprintf("└┬─> %s", sfyTask["name"])
				}

				printStringyTask(w, sfyTask)

				leaves := make([]uint64, len(sub.nodes))
				j := 0
				for k := range sub.nodes {
					leaves[j] = k
					j++
				}
				sort.Slice(leaves, func(i, j int) bool { return leaves[i] < leaves[j] })

				j = 0
				for _, lkey := range leaves {
					leaf := sub.nodes[lkey]
					if len(leaf.nodes) != 0 {
						// TODO: Supersubs
					}

					leafTask, err := get(leaf.key)
					// TODO: Log
					if err != nil {
						fmt.Println("Could not get task %d: %s\n", leaf.key, err.Error())
					}

					sfyTask = leafTask.stringify()

					if len(sub.nodes) > 1 && len(sub.nodes) != j+1 {
						sfyTask["name"] = fmt.Sprintf("%s├──> %s", prefix, sfyTask["name"])
					} else {
						sfyTask["name"] = fmt.Sprintf("%s└──> %s", prefix, sfyTask["name"])
					}

					printStringyTask(w, sfyTask)
					j++
				}
			} else {
				if len(root.nodes) > 1 && len(root.nodes) != i+1 {
					sfyTask["name"] = fmt.Sprintf("├──> %s", sfyTask["name"])
				} else {
					sfyTask["name"] = fmt.Sprintf("└──> %s", sfyTask["name"])
				}

				printStringyTask(w, sfyTask)
			}

			i++
		}
	}

	w.Flush()
}
