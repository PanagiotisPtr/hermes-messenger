package user

import "github.com/google/uuid"

// User entity used to store information about the User
type User struct {
	ID uuid.UUID `bson:"_id" json:"id"`
	UserDetails
}

// UserDetails struct representing the details available
// for a user
type UserDetails struct {
	Email     string `bson:"email" json:"email"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
}
