package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/maxcelant/istio-microservice-sample-items/internal/svc"
)

func initDB() (*sql.DB, func()) {
	connStr := "postgres://items_user:items_pass@localhost:5433/items_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	cleanup := func() {
		db.Close()
	}

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		defer cleanup()
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		defer cleanup()
	}

	return db, cleanup
}

func main() {
	db, cleanup := initDB()
	defer cleanup()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
		defer cleanup()
	}

	lg := log.New(os.Stdout, "items ", log.LstdFlags)
	itemService := svc.New(db, lg)
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/api/items", itemService.GetItems()).Methods(http.MethodGet)
	router.HandleFunc("/api/items", itemService.CreateItem()).Methods(http.MethodPost)
	router.HandleFunc("/api/items/{id}", itemService.GetItemByID()).Methods(http.MethodGet)

	s := &http.Server{
		Addr:         ":8081",
		Handler:      router,
		ErrorLog:     lg,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	lg.Print("Listening on port 8081")
	lg.Fatal(s.ListenAndServe())
}
