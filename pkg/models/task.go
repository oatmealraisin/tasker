package models

import (
	"math"
	"sort"
	"time"

	"github.com/golang/protobuf/ptypes"
)

// High score means more important
func (task Task) Score() float64 {
	add_time, err := ptypes.Timestamp(task.Added)
	if err != nil {
		return 0.0
	}

	due_mod := 0.0
	if task.Due != nil {
		due_time, err := ptypes.Timestamp(task.Due)
		if err != nil {
			return 0.0
		}

		due_mod = 24.0 / math.Exp(time.Until(due_time).Hours()*0.05)
	}

	add_mod := math.Max(1.0, math.Log(time.Since(add_time).Hours()/24.0))

	// TODO: Size and priority have a special relationship.. You want to do the
	// smallest, most important tasks first, followed by the hardest, most
	// important tasks. Change the formula to reflect this
	// TODO: Also, the age kind of changes the priority.. or at least makes the
	// priority less important. Should also change the formula to reflect this
	priority_mod := math.Pow(3.0, float64(task.Priority)) * 0.075
	size_mod := 0.5 * float64(task.Size)

	return due_mod*size_mod + add_mod/(priority_mod+size_mod)
}

func UuidSort(uuids []uint64) {
	if len(uuids) < 2 {
		return
	}

	if len(uuids) < 256 {
		sort.Slice(uuids, func(i, j int) bool { return uuids[i] < uuids[j] })
	}

	buffer := make([]uint64, len(uuids))

	// Each pass processes a byte offset, copying back and forth between slices
	from := uuids
	to := buffer[:len(uuids)]
	var key uint8
	var offset [256]int // Keep track of where groups start

	for keyOffset := uint(0); keyOffset < 64; keyOffset += 8 {
		keyMask := uint64(0xFF << keyOffset) // Current 'digit' to look at
		var counts [256]int                  // Keep track of the number of elements for each kind of byte
		sorted := true                       // Check for already sorted
		prev := uint64(0)                    // if elem is always >= prev it is already sorted
		for _, elem := range from {
			key = uint8((elem & keyMask) >> keyOffset) // fetch the byte at current 'digit'
			counts[key]++                              // count of elems to put in this digit's bucket

			if sorted { // Detect sorted
				sorted = elem >= prev
				prev = elem
			}
		}

		if sorted { // Short-circuit sorted
			if (keyOffset/8)%2 == 1 {
				copy(to, from)
			}
			return
		}

		// Find target bucket offsets
		offset[0] = 0
		for i := 1; i < len(offset); i++ {
			offset[i] = offset[i-1] + counts[i-1]
		}

		// Rebucket while copying to other buffer
		for _, elem := range from {
			key = uint8((elem & keyMask) >> keyOffset) // Get the digit
			to[offset[key]] = elem                     // Copy the element to the digit's bucket
			offset[key]++                              // One less space, move the offset
		}
		// On next pass copy data the other way
		to, from = from, to
	}
}

func GetAllChildren(uuid uint64, get func(t uint64) (Task, error)) ([]uint64, error) {
	result := []uint64{uuid}
	var err error

	task, err := get(uuid)
	if err != nil {
		return result, err
	}

	for _, child := range task.Subtasks {
		c, err := GetAllChildren(child, get)
		if err != nil {
			return []uint64{}, err
		}

		result = append(result, c...)
	}

	return result, err
}
