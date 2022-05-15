package config

import "log"

type Service struct {
	logger    *log.Logger
	paramRepo Repository
}

func NewService(logger *log.Logger) *Service {
	return &Service{
		logger:    logger,
		paramRepo: NewMemoryRepository(logger),
	}
}

func (s *Service) SetParam(name string, value string) error {
	return s.paramRepo.AddParam(name, value)
}

func (s *Service) UnsetParam(name string) error {
	return s.paramRepo.RemoveParam(name)
}

func (s *Service) GetParam(name string) (*Param, error) {
	return s.paramRepo.GetParam(name)
}
