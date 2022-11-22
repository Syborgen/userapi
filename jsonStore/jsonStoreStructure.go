package jsonstore

import "refactoring/store"

type jsonStoreStructure struct {
	Increment int                   `json:"increment"`
	Users     map[string]store.User `json:"users"`
}
