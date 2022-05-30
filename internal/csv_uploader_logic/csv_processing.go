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
	tempFile, _ = os.Open(tempFile.Name())
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	fcsv := csv.NewReader(tempFile)
	numWps := 100
	jobs := make(chan []string, numWps)
	res := make(chan *api.Order)

	orders := make([]*api.Order, 0)

	var wg sync.WaitGroup
	worker := func(jobs <-chan []string, results chan<- *api.Order) {
		for {
			select {
			case job, ok := <-jobs: // Check for readable state of the channel.
				if !ok {
					return
				}
				result, err := parseStruct(job) //Parse the line into a struct
				if err != nil {
					fmt.Println(err)
					return
				}
				// Send the object to the channel after processed
				results <- result
			}
		}
	}

	// init workers
	for w := 0; w < numWps; w++ {
		wg.Add(1)
		go func() {
			// this line will exec when chan `res` processed output
			defer wg.Done()
			worker(jobs, res)
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
		wg.Wait()
		close(res) // when you close(res) it breaks the below loop.
	}()

	for r := range res {
		orders = append(orders, r) // Append all processed files into a slice
	}

	// Call api to save the data into a database
	errPost := callPersistenceApi(orders)

	if errPost != nil {
		fmt.Println(errPost)
		return
	}

	fmt.Println("Processed ", len(orders))
}

func callPersistenceApi(orders []*api.Order) error {
	jsonData, err := json.Marshal(orders) // Change the sclice into a json array
	if err != nil {
		return err
	}
	// Do the POST Request
	response, err := http.Post("http://localhost:9001/orders", "application/json", bytes.NewBuffer(jsonData))

	// Return generic error
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Check StatusCode for Created
	if response.StatusCode != 201 {
		err := errors.New("couldn't create resource") //New error for coudn't create
		return err
	}

	return nil
}

func parseStruct(data []string) (*api.Order, error) {
	id, _ := strconv.ParseInt(strings.TrimSpace(data[0]), 10, 32)
	email := strings.TrimSpace(data[1])
	phoneNumber := phoneNumberFormat(data[2])

	parcelWeight, err := strconv.ParseFloat(strings.TrimSpace(data[3]), 32)
	if err != nil {
		return nil, err
	}

	year, month, day := time.Now().Date()
	date := strconv.Itoa(year) + "-" + month.String() + "-" + strconv.Itoa(day)

	country, err := getCountry(phoneNumber)

	if err != nil {
		return nil, err
	}

	return api.CreateOrder(int(id), email, phoneNumber, float32(parcelWeight), date, country), nil
}

func phoneNumberFormat(phoneNumber string) string {
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "") // Clear all spaces, was not needed. Just a nice touch
	return "(" + phoneNumber[0:3] + ")" + phoneNumber[3:]  // Create a format to regex recognize
}

func getCountry(phoneNumber string) (string, error) {
	for country, reg := range countriesRegex { // Cycle through regex list
		if reg.MatchString(phoneNumber) { // When Match
			return country, nil // Return the Country
		}
	}

	// If couldn't find a match, return a error
	return "", errors.New("country not found")
}
