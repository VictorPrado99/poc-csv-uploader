package controller

import "net/http"

func SetupRoutes() {
	http.HandleFunc("/upload", uploadCsv)
}

func uploadCsv(w http.ResponseWriter, r *http.Request) {

}
