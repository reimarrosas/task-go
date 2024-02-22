package models

import (
	"strconv"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("task-bucket")

type Todo struct {
	ID          []byte
	Description string
}

type Task struct {
	db *bolt.DB
}

func InitTask(dbName string) (*Task, error) {
	db, err := bolt.Open(dbName, 0600, nil)

	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket(taskBucket)
		return nil
	})

	return &Task{db}, nil
}

func (t *Task) Add(todo string) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		id64, err := b.NextSequence()
		if err != nil {
			return err
		}

		id := itoa(id64)
		if err := b.Put(id, []byte(todo)); err != nil {
			return err
		}

		return nil
	})
}

func (t *Task) Delete(id string) error {
	return t.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		ts, err := t.List()
		if err != nil {
			return err
		}

		idI, err := atoi(id)
		if err != nil {
			return err
		}

		for i, todo := range ts {
			if i == idI-1 {
				if err := b.Delete(todo.ID); err != nil {
					return nil
				}
			}
		}

		return nil
	})
}

func (t *Task) List() ([]Todo, error) {
	var res []Todo
	err := t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			res = append(res, Todo{k, string(v)})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func itoa[T uint64 | int64 | int32 | int](i T) []byte {
	a := strconv.Itoa(int(i))
	return []byte(a)
}

func atoi(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return i, nil
}
