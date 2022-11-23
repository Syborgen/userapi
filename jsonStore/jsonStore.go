package jsonstore

import (
	"encoding/json"
	"fmt"
	"os"
	"refactoring/store"
	"strconv"
)

type JsonStore struct {
	FileName string

	lastUserIndex int
}

func (js *JsonStore) AddUser(user store.User) (int, error) {
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

func (js *JsonStore) GetUsers() (store.Users, error) {
	return js.readUsers()
}

func (js *JsonStore) GetUser(id string) (store.User, error) {
	users, err := js.readUsers()
	if err != nil {
		return store.User{}, fmt.Errorf("read users error: %w", err)
	}

	user, ok := users[id]
	if !ok {
		return store.User{}, fmt.Errorf("user with id %s is not exists", id)
	}

	return user, nil
}

func (js *JsonStore) readUsers() (store.Users, error) {
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

func (js *JsonStore) writeUsers(users map[string]store.User) error {
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
