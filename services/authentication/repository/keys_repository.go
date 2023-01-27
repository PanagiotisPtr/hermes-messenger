package repository

import "context"

type KeysRepository interface {
	GenerateKeys(ctx context.Context) error
	GetPublicKeys(ctx context.Context) ([]string, error)
}
