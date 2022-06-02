package csvuploaderlogic

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	api "github.com/VictorPrado99/poc-csv-persistence/pkg/api"
)

func ProcessCsv(tempFile *os.File) {

	orders := DoCsvProcess(tempFile)

	// Call api to save the data into a database
	errPost := callPersistenceApi(orders)

	// If got some error during POST
	if errPost != nil {
		fmt.Println(errPost)
		return
	}

	// Log how many objects was processed
	fmt.Println("Processed ", len(orders))
}

func DoCsvProcess(tempFile *os.File) []*api.Order {
	tempFile, _ = os.Open(tempFile.Name()) //Open the temp file
	defer tempFile.Close()                 // Close when finish
	defer os.Remove(tempFile.Name())       // Remove when finis

	fcsv := csv.NewReader(tempFile)     //Start the csv reader
	numWps := 100                       // Number of worker, can be higher if you want more paralellism
	jobs := make(chan []string, numWps) //A channel to send lines to be processed
	res := make(chan *api.Order)        //Cahnnel to receive the processed result

	orders := make([]*api.Order, 0)

	var wg sync.WaitGroup // WaitGroup to control the jobs
	worker := func(jobs <-chan []string, results chan<- *api.Order) {
		for {
			select {
			case job, ok := <-jobs: // Check for readable state of the channel.
				if !ok {
					return
				}
				result, err := parseStruct(job) //Parse the line into a struct
				if err != nil {                 // If coudn't parse it, ignore it... Could be stored to analysed later, or maybe show the client what went wrong
					fmt.Println(err)
					return
				}
				// Send the object to the channel after processed
				results <- result
			}
		}
	}

	// init workers
	for w := 0; w < numWps; w++ { // for each worker
		wg.Add(1) // Add worker to counter
		go func() {
			defer wg.Done()   // When worker finish his job, tell the wait group we are done
			worker(jobs, res) // Instantiate worker
		}()
	}

	go func() {
		firstLine := true
		for {
			rStr, err := fcsv.Read()

			// Ignore the first line, we don't need to process the header
			if firstLine {
				firstLine = false
				continue
			}

			// Break the loop at the end of file
			if err == io.EOF {
				break
			}

			// If something has gone wrong stop the process
			if err != nil {
				fmt.Println("ERROR: ", err.Error())
				break
			}

			// Ready the line to be processed by the worker
			jobs <- rStr
		}
		close(jobs) // close jobs to signal workers that no more job are incoming.
	}()

	go func() {
		wg.Wait()  // When all workers finish their jobs
		close(res) // close(res) it breaks the below loop.
	}()

	for r := range res {
		orders = append(orders, r) // Append all processed files into a slice
	}

	return orders
}

func callPersistenceApi(orders []*api.Order) error {
	jsonData, err := json.Marshal(orders) // Change the sclice into a json array
	if err != nil {
		return err
	}

	// Get envrioment variable. This way, supports hot deploy.
	persistenceUrl := os.Getenv("PERSIST_URL")

	if persistenceUrl == "" {
		persistenceUrl = "http://localhost:9001" // Default value
	}

	// Do the POST Request
	response, err := http.Post(persistenceUrl+"/orders", "application/json", bytes.NewBuffer(jsonData))

	// Return generic error
	if err != nil {
		fmt.Println(err) // Could log or send for some observability platform
		return err
	}

	// Check StatusCode for Created
	if response.StatusCode != 201 {
		err := errors.New("couldn't create resource") //New error for couldn't create
		return err
	}

	return nil
}

func parseStruct(data []string) (*api.Order, error) {
	// Get fields
	id, _ := strconv.ParseInt(strings.TrimSpace(data[0]), 10, 32)
	email := strings.TrimSpace(data[1])
	phoneNumber := phoneNumberFormat(data[2])

	parcelWeight, err := strconv.ParseFloat(strings.TrimSpace(data[3]), 32)
	if err != nil { //Could't parse float values
		return nil, err
	}

	date := getDate()

	// Get country based on phone number cycling through a regex dictionary
	country, err := getCountry(phoneNumber)

	// If regex failed, or couldn't find a country
	if err != nil {
		return nil, err
	}

	// Otherwise return the built object
	return api.CreateOrder(int(id), email, phoneNumber, float32(parcelWeight), date, country), nil
}

// Get date base upon this exacly moment
func getDate() string {
	year, month, day := time.Now().Date() // Get date
	monthStr := strconv.Itoa(int(month))
	date := strconv.Itoa(year) + "-" + monthStr + "-" + strconv.Itoa(day) // Build date at my format
	return date
}

//Format the phone number with (XXX)XXXXXX... For aesthetic and regex purposes
func phoneNumberFormat(phoneNumber string) string {
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "") // Clear all spaces, was not needed. Just a nice touch
	return "(" + phoneNumber[0:3] + ")" + phoneNumber[3:]  // Create a format to regex recognize
}

// Get country based upon the phone number. Use regex to determine it.
func getCountry(phoneNumber string) (string, error) {
	for country, reg := range countriesRegex { // Cycle through regex list
		if reg.MatchString(phoneNumber) { // When Match
			return country, nil // Return the Country
		}
	}

	// If couldn't find a match, return a error
	return "", errors.New("country not found")
}
