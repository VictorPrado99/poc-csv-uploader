package main

import (
	"log"
	"net/http"
	"time"

	"github.com/VictorPrado99/poc-csv-uploader/internal/controller"
)

func main() {
	router := controller.SetupRoutes()
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:9100",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
