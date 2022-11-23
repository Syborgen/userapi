package jsonstore

import (
	"encoding/json"
	"errors"
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
	if !js.isFileExists() {
		err := js.initStorageFile()
		if err != nil {
			return nil, fmt.Errorf("init json storage file error: %w", err)
		}
	}

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

func (js *JSONStore) isFileExists() bool {
	_, err := os.Stat(js.FileName)
	return !errors.Is(err, os.ErrNotExist)
}

func (js *JSONStore) initStorageFile() error {
	file, err := os.Create(js.FileName)
	if err != nil {
		return fmt.Errorf("create storage file '%s' error: %w",
			js.FileName, err)
	}

	defer file.Close()

	initialStorageStructure, err := json.Marshal(map[string]interface{}{
		"increment": 0,
		"users":     struct{}{},
	})
	if err != nil {
		return fmt.Errorf("json marshaling error: %w", err)
	}

	_, err = file.Write(initialStorageStructure)
	if err != nil {
		return fmt.Errorf("write to file error: %w", err)
	}

	return nil
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
