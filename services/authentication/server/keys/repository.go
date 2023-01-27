package keys

import (
	"context"
)

type Repository interface {
	GenerateKeys(context.Context) error
}
