package handlers

import (
	"net/http"
	"refactoring/helper"
	jsonstore "refactoring/jsonStore"
	"refactoring/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const storeFileName = `users.json`

var dataStorage store.Store = &jsonstore.JSONStore{FileName: storeFileName}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := dataStorage.GetUsers()
	if err != nil {
		render.Render(w, r, helper.ErrInternalError(err))
		return
	}

	render.JSON(w, r, users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	request := CreateUserRequest{}
	err := render.Bind(r, &request)
	if err != nil {
		render.Render(w, r, helper.ErrInvalidRequest(err))
		return
	}

	newUser := store.User{
		CreatedAt:   time.Now(),
		DisplayName: request.DisplayName,
		Email:       request.Email,
	}
	newUserID, err := dataStorage.AddUser(newUser)
	if err != nil {
		render.Render(w, r, helper.ErrInternalError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": newUserID,
	})
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := dataStorage.GetUser(id)
	if err != nil {
		render.Render(w, r, helper.ErrInternalError(err))
		return
	}

	render.JSON(w, r, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	request := CreateUserRequest{}
	err := render.Bind(r, &request)
	if err != nil {
		render.Render(w, r, helper.ErrInvalidRequest(err))
		return
	}

	newUserData := store.User{
		DisplayName: request.DisplayName,
		Email:       request.Email,
	}
	err = dataStorage.UpdateUser(id, newUserData)
	if err != nil {
		render.Render(w, r, helper.ErrInternalError(err))
		return
	}

	render.NoContent(w, r)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := dataStorage.DeleteUser(id)
	if err != nil {
		render.Render(w, r, helper.ErrInternalError(err))
		return
	}

	render.NoContent(w, r)
}
