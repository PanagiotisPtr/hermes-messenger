package credentials

import "github.com/google/uuid"

type Credentials struct {
	Uuid     uuid.UUID `json:"Uuid"`
	Password string    `json:"Password"`
}
