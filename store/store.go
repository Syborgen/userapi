package store

type Store interface {
	AddUser(user User) (int, error)
	GetUsers() (Users, error)
	GetUser(id string) (User, error)
	DeleteUser(id string) error
	UpdateUser(id string, newUserData User) error
}
