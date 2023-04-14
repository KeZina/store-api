package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {
	var conn string
	var port string

	flag.StringVar(&conn, "conn", "", "db connection string")
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

	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Poof!!!"))
		if err != nil {
			log.Println(err)
		}
	})

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Println(err)
	}
}
