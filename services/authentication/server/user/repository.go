package user

type Repository interface {
	CreateUser(user User) error
	GetUserByEmail(email string) (User, error)
	DeleteUser(user User) error
}
