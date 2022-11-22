package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	jsonstore "refactoring/jsonStore"
	"refactoring/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const storeFileName = `users.json`

var dataStorage store.Store = &jsonstore.JsonStore{FileName: storeFileName}

type (
	UserStore struct {
		Increment int         `json:"increment"`
		List      store.Users `json:"list"`
	}
)

var (
	UserNotFound = errors.New("user_not_found")
)

func main() {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", searchUsers)
				r.Post("/", createUser)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", getUser)
					r.Patch("/", updateUser)
					r.Delete("/", deleteUser)
				})
			})
		})
	})

	fmt.Println("Start server.")
	http.ListenAndServe(":3333", router)
}

func searchUsers(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(storeFileName)
	s := UserStore{}
	_ = json.Unmarshal(f, &s)

	render.JSON(w, r, s.List)
}

type CreateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (c *CreateUserRequest) Bind(r *http.Request) error {
	return nil
}

func createUser(w http.ResponseWriter, r *http.Request) {
	request := CreateUserRequest{}
	err := render.Bind(r, &request)
	if err != nil {
		err := render.Render(w, r, ErrInvalidRequest(err))
		if err != nil {
			fmt.Println("Render error:", err)
		}

		return
	}

	newUser := store.User{
		CreatedAt:   time.Now(),
		DisplayName: request.DisplayName,
		Email:       request.Email,
	}
	newUserID, err := dataStorage.AddUser(newUser)
	if err != nil {
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": newUserID,
	})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(storeFileName)
	s := UserStore{}
	_ = json.Unmarshal(f, &s)

	id := chi.URLParam(r, "id")

	render.JSON(w, r, s.List[id])
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
}

func (c *UpdateUserRequest) Bind(r *http.Request) error { return nil }

func updateUser(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(storeFileName)
	s := UserStore{}
	_ = json.Unmarshal(f, &s)

	request := UpdateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	id := chi.URLParam(r, "id")

	if _, ok := s.List[id]; !ok {
		_ = render.Render(w, r, ErrInvalidRequest(UserNotFound))
		return
	}

	u := s.List[id]
	u.DisplayName = request.DisplayName
	s.List[id] = u

	b, _ := json.Marshal(&s)
	_ = ioutil.WriteFile(storeFileName, b, fs.ModePerm)

	render.Status(r, http.StatusNoContent)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	f, _ := ioutil.ReadFile(storeFileName)
	s := UserStore{}
	_ = json.Unmarshal(f, &s)

	id := chi.URLParam(r, "id")

	if _, ok := s.List[id]; !ok {
		_ = render.Render(w, r, ErrInvalidRequest(UserNotFound))
		return
	}

	delete(s.List, id)

	b, _ := json.Marshal(&s)
	_ = ioutil.WriteFile(storeFileName, b, fs.ModePerm)

	render.Status(r, http.StatusNoContent)
}

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
