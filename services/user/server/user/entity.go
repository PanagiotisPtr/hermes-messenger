package user

import "github.com/google/uuid"

type User struct {
	// Could store more user details here like first/last names etc.
	Uuid  uuid.UUID `bson:"_id" json:"Uuid"`
	Email string    `bson:"Email" json:"Email"`
}
