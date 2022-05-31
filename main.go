package main

import (
	"os"

	service "github.com/VictorPrado99/poc-csv-uploader/cmd/csv_uploader_service"
)

func main() {
	port := os.Getenv("CSV_UPLOAD_PORT")

	if port == "" {
		port = "9100"
	}

	service.ServiceStart(port)
}
