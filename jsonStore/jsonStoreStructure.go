package jsonstore

import "refactoring/store"

type jsonStoreStructure struct {
	Increment int         `json:"increment"`
	Users     store.Users `json:"users"`
}
