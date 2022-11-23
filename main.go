package main

import (
	"fmt"
	"net/http"
	"refactoring/handlers"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", handlers.GetUsers)
				r.Post("/", handlers.CreateUser)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", handlers.GetUser)
					r.Patch("/", handlers.UpdateUser)
					r.Delete("/", handlers.DeleteUser)
				})
			})
		})
	})

	fmt.Println("Start server.")
	http.ListenAndServe(":3333", router)
}
