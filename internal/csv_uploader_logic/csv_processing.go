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
				result, err := parseStruct(job)
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
		for {
			rStr, err := fcsv.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("ERROR: ", err.Error())
				break
			}
			jobs <- rStr
		}
		close(jobs) // close jobs to signal workers that no more job are incoming.
	}()

	go func() {
		wg.Wait()
		close(res) // when you close(res) it breaks the below loop.
	}()

	for r := range res {
		orders = append(orders, r)
	}

	callPersistenceApi(orders)
	fmt.Println("Processed ", len(orders))
}

func callPersistenceApi(orders []*api.Order) error {
	jsonData, err := json.Marshal(orders)
	if err != nil {
		return err
	}
	response, err := http.Post("http://localhost:9001/orders", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return err
	}

	if response.StatusCode != 201 {
		err := errors.New("couldn't create resource")
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
	phoneNumber = strings.ReplaceAll(strings.TrimSpace(phoneNumber), " ", "")
	return "(" + phoneNumber[0:3] + ")" + phoneNumber[3:]
}

func getCountry(phoneNumber string) (string, error) {
	for country, reg := range countriesRegex {
		if reg.MatchString(phoneNumber) {
			return country, nil
		}
	}
	return "", errors.New("country not found")
}
