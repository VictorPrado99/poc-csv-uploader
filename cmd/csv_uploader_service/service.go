package csvuploaderservice

import (
	"log"
	"net/http"
	"time"

	controller "github.com/VictorPrado99/poc-csv-uploader/internal/csv_uploader_controller"
)

func ServiceStart() {
	router := controller.SetupRoutes()
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:9100",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
