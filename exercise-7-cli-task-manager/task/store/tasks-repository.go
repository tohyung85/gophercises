package store

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
)

type BoltStore struct {
	DB *bolt.DB
}

type Task struct {
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func NewStore(db *bolt.DB) *BoltStore {
	return &BoltStore{db}
}

func (bs *BoltStore) AddToStore(item string) error {
	return bs.DB.Update(func(tx *bolt.Tx) error {
		task := Task{item, false}
		taskBytes, err := json.Marshal(task)
		if err != nil {
			return err
		}
		b := tx.Bucket([]byte("DB"))
		id, _ := b.NextSequence()
		e := b.Put(itob(int(id)), taskBytes)
		return e
	})
}

func (bs *BoltStore) RetrieveAll() (map[int]Task, error) {
	allTasks := make(map[int]Task)
	e := bs.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var task Task
			err := json.Unmarshal(v, &task)
			if err != nil {
				return err
			}
			allTasks[btoi(k)] = task
		}
		return nil
	})
	return allTasks, e
}

func (bs *BoltStore) FlagComplete(itemNo int) (Task, error) {
	var task Task
	e := bs.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB"))
		_, v := b.Cursor().Seek(itob(itemNo))
		err := json.Unmarshal(v, &task)
		if err != nil {
			return err
		}
		task.Completed = true
		taskBytes, err := json.Marshal(task)
		if err != nil {
			return err
		}
		return b.Put(itob(itemNo), taskBytes)
	})
	return task, e
}

func (bs *BoltStore) DeleteItem(itemNo int) error {
	e := bs.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB"))
		err := b.Delete(itob(itemNo))
		if err != nil {
			return err
		}
		return nil
	})
	return e
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
