package csvuploadercontroller

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	logic "github.com/VictorPrado99/poc-csv-uploader/internal/csv_uploader_logic"

	"github.com/gorilla/mux"
)

const (
	GET  = "GET"
	POST = "POST"
	PUT  = "PUT"

	TEMP_DIR = "./temp-files"
)

// Method who will setup the router of this controller
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/upload", UploadCsv).Methods(POST)

	return r
}

// Check if the first line of csv is as expected
func checkCsvHeader(header []string) bool {
	return strings.TrimSpace(header[0]) == "id" && strings.TrimSpace(header[1]) == "email" && strings.TrimSpace(header[2]) == "phone_number" && strings.TrimSpace(header[3]) == "parcel_weight"
}

// Method to receive a CSV and made a async processing
func UploadCsv(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(0) //Receive csv
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(r.Header.Get("content-type"))

	file, _, err := r.FormFile("myCsv")
	if err != nil { // Coudn't retrieve the file
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	// Read the first csv line
	fcsv := csv.NewReader(file)
	line, err := fcsv.Read()

	if err != nil {
		//Return not accepted
		w.WriteHeader(http.StatusNotAcceptable)

		// return that the error to the client
		fmt.Fprintf(w, "Couldn't read File, %v\n", err)
		return
	}

	// Check the header follow the contract
	if !checkCsvHeader(line) { //if not
		//Return not accepted
		w.WriteHeader(http.StatusNotAcceptable)

		// return that we have successfully uploaded our file!
		fmt.Fprintf(w, "Header do not conform with the interface\n")
		return
	}

	// Close the file when the function finish it
	defer file.Close()

	// Create a directory if needed, this temp dir could eventually be a shared folder for the pod. Guided by events to be more distributed
	os.MkdirAll(TEMP_DIR, os.ModePerm)

	// Create a temporary file within our temp-file directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile(TEMP_DIR, "upload-*.csv")
	if err != nil {
		fmt.Println(err)
	}
	// Close the file when the function end
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	// Start the processing in another goroutine, to keep it simple
	go logic.ProcessCsv(tempFile)

	//Return accepted, and process the file in background
	w.WriteHeader(http.StatusAccepted)

	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
