package store

type Store interface {
	AddUser(user User) (int, error)
	GetUsers() (Users, error)
}
