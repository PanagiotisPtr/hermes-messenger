package config

type Repository interface {
	AddParam(string, string) error
	RemoveParam(string) error
	GetParam(string) (*Param, error)
}
