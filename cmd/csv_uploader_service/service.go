package csvuploaderservice

import (
	"fmt"
	"log"
	"net/http"
	"time"

	controller "github.com/VictorPrado99/poc-csv-uploader/internal/csv_uploader_controller"
)

func ServiceStart(port string) {
	router := controller.SetupRoutes()
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf(fmt.Sprintf("Starting Server on port %s \n", port))
	log.Fatal(srv.ListenAndServe())
}
