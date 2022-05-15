package config

import (
	"fmt"
	"log"
)

type MemoryRepository struct {
	logger *log.Logger
	params map[string]*Param
}

func NewMemoryRepository(logger *log.Logger) *MemoryRepository {
	return &MemoryRepository{
		logger: logger,
		params: make(map[string]*Param),
	}
}

func (mr *MemoryRepository) AddParam(name string, value string) error {
	if _, ok := mr.params[name]; ok {
		return fmt.Errorf("Parameter with name '%s' already exists", name)
	}

	mr.params[name] = &Param{
		Name:  name,
		Value: value,
	}

	return nil
}

func (mr *MemoryRepository) RemoveParam(name string) error {
	if _, ok := mr.params[name]; !ok {
		return nil
	}
	delete(mr.params, name)

	return nil
}

func (mr *MemoryRepository) GetParam(name string) (*Param, error) {
	param, ok := mr.params[name]
	if !ok {
		return nil, fmt.Errorf("Parameter with name '%s' does not exist", name)
	}

	return param, nil
}
