package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"proj/internal/store"
	"proj/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

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

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Recoverer)

	userRepo := user.UserRepository{DB: db}
	storeRepo := store.StoreRepository{DB: db}

	userService := user.UserService{UserRepo: userRepo}
	storeService := store.StoreService{StoreRepo: storeRepo}

	r.Group(func(r chi.Router) {
		r.Use(user.Auth)

		r.Get("/users/profile", userService.GetUser)

		r.Get("/store/users/items", storeService.GetUserStoreItems)
		r.Delete("/store/users/items", storeService.DeleteUserStoreItem)
		r.Get("/store/available-items", storeService.GetAvailableStoreItems)
		r.Post("/store/purchase", storeService.PurchaseStoreItem)
	})

	r.Group(func(r chi.Router) {
		r.Post("/users/login", userService.Login)
		r.Get("/users/logout", userService.Logout)
	})

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Println(err)
	}
}
