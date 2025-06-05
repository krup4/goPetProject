package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/krup4/goPetProject/web/backend/db"
	"github.com/krup4/goPetProject/web/backend/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := db.ConnectToDB(os.Getenv("POSTGRES_CONN"))
	if err != nil {
		log.Fatal(err)
	}

	err = db.InitDB(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	r.Get("/api/v1/ping", handlers.PingHandler)

	r.Route("/api/v1/user", func(r chi.Router) {
		r.Post("/sign-up", handlers.UserSignUP)
	})

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
