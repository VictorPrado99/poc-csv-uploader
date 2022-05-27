package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	GET  = "GET"
	POST = "POST"
	PUT  = "PUT"
)

// Method who will setup the router of this controller
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/upload", uploadCsv).Methods(POST)

	return r
}

// Method to receive a CSV and made a async processing
func uploadCsv(w http.ResponseWriter, r *http.Request) {
	fmt.Println("upload method called")
	err := r.ParseMultipartForm(0)
	if err != nil {
		fmt.Println(err)
	}

	file, handler, err := r.FormFile("myCsv")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp-files", "upload-*.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
