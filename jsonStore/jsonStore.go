package jsonstore

import (
	"encoding/json"
	"fmt"
	"os"
	"refactoring/helper"
	"refactoring/store"
	"strconv"
)

type JSONStore struct {
	FileName string

	lastUserIndex int
}

func (js *JSONStore) AddUser(user store.User) (int, error) {
	users, err := js.readUsers()
	if err != nil {
		return 0, fmt.Errorf("read users error: %w", err)
	}

	js.lastUserIndex++
	users[strconv.Itoa(js.lastUserIndex)] = user

	err = js.writeUsers(users)
	if err != nil {
		return 0, fmt.Errorf("write users error: %w", err)
	}

	return js.lastUserIndex, nil
}

func (js *JSONStore) GetUsers() (store.Users, error) {
	return js.readUsers()
}

func (js *JSONStore) GetUser(id string) (store.User, error) {
	users, err := js.readUsers()
	if err != nil {
		return store.User{}, fmt.Errorf("read users error: %w", err)
	}

	user, ok := users[id]
	if !ok {
		return store.User{}, helper.ErrUserNotFound
	}

	return user, nil
}

func (js *JSONStore) DeleteUser(id string) error {
	users, err := js.readUsers()
	if err != nil {
		return fmt.Errorf("read users error: %w", err)
	}

	if _, ok := users[id]; !ok {
		return helper.ErrUserNotFound
	}

	delete(users, id)

	err = js.writeUsers(users)
	if err != nil {
		return fmt.Errorf("write users error: %w", err)
	}

	return nil
}

func (js *JSONStore) UpdateUser(id string, newUserData store.User) error {
	users, err := js.readUsers()
	if err != nil {
		return fmt.Errorf("read users error: %w", err)
	}

	userToUpdate, ok := users[id]
	if !ok {
		return helper.ErrUserNotFound
	}

	userToUpdate.Update(newUserData)
	users[id] = userToUpdate

	err = js.writeUsers(users)
	if err != nil {
		return fmt.Errorf("read users error: %w", err)
	}

	return nil
}

func (js *JSONStore) readUsers() (store.Users, error) {
	fileData, err := os.ReadFile(js.FileName)
	if err != nil {
		return nil, fmt.Errorf("read file error: %w", err)
	}

	var jsonData jsonStoreStructure
	err = json.Unmarshal(fileData, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	js.lastUserIndex = jsonData.Increment
	return jsonData.Users, nil
}

func (js *JSONStore) writeUsers(users map[string]store.User) error {
	dataToWrite := jsonStoreStructure{
		Increment: js.lastUserIndex,
		Users:     users,
	}

	jsonData, err := json.Marshal(&dataToWrite)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	err = os.WriteFile(js.FileName, jsonData, os.ModePerm)
	if err != nil {
		return fmt.Errorf("write to file error: %w", err)
	}

	return nil
}
