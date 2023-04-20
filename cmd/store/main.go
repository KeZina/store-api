package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"proj/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	var conn string
	var port string

	flag.StringVar(&conn, "conn", "host=localhost port=5432 user=postgres password=b00mka dbname=test sslmode=disable", "db connection string")
	flag.StringVar(&port, "port", "8080", "server port")
	flag.Parse()

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	userRepo := user.UserRepository{DB: db}
	userService := user.UserService{UserRepo: userRepo}

	r.Group(func(r chi.Router) {
		r.Use(user.Auth)

		r.Post("/user/ping", userService.Ping)
	})

	r.Group(func(r chi.Router) {
		r.Post("/user/login", userService.Login)
	})

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Println(err)
	}
}
