package main

import (
	"fmt"
	"go_server/internal/handlers"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func main() {
	// r := chi.NewRouter()           // Use chi
	// log.Info("Starting server...") // Use logrus
	// fmt.Println("Hello, world!")   // Use fmt
	// http.Handle("/", r)            // Use net/http
	// // Example: handlers.SomeFunction() // Use your internal package
	// http.ListenAndServe(":8080", nil)

	log.SetReportCaller(true)
	var r *chi.Mux = chi.NewRouter()

	handlers.Handler(r)
	fmt.Println("Starting GO API Services...")
	fmt.Println(`go go go api`)

	err := http.ListenAndServe("Localhost:8000", r)

	if err != nil {
		log.Error(err)
	}
}
