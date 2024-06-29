package storage

import (
	"encoding/json"
	"errors"
	"os"
)

var (
	ErrNotFound = errors.New("not found")
	ErrSave     = errors.New("error saving data")
)

type Storage struct {
	storeFunc func(key string, value json.RawMessage) error
	getFunc   func(key string) (json.RawMessage, error)
}

type ConfigFunc func(*Storage)

func New(dataFile string) (*Storage, error) {
	s := &Storage{}

	values, err := loadData(dataFile)
	if err != nil {
		return nil, err
	}

	s.storeFunc = func(key string, value json.RawMessage) error {
		// make copy of map
		newValues := make(map[string]json.RawMessage)
		for k, v := range values {
			newValues[k] = v
		}

		newValues[key] = value

		err := saveData(dataFile, newValues)
		if err != nil {
			return errors.Join(ErrSave, err)
		}

		values = newValues
		return nil
	}

	s.getFunc = func(key string) (json.RawMessage, error) {
		v, ok := values[key]
		if !ok {
			return nil, ErrNotFound
		}
		return v, nil
	}

	return s, nil
}

func (s *Storage) Store(key string, value json.RawMessage) error {
	return s.storeFunc(key, value)
}

func (s *Storage) Get(key string) (json.RawMessage, error) {
	return s.getFunc(key)
}

func loadData(file string) (map[string]json.RawMessage, error) {
	var values map[string]json.RawMessage

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func saveData(file string, values map[string]json.RawMessage) error {
	data, err := json.MarshalIndent(values, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(file, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
