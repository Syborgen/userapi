package store

type Store interface {
	AddUser(user User) (int, error)
}
