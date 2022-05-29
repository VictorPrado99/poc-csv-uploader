package logic

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

	"github.com/VictorPrado99/poc-csv-uploader/pkg/api"
)

func ProcessCsv(tempFile *os.File) {
	tempFile, _ = os.Open(tempFile.Name())
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	fcsv := csv.NewReader(tempFile)
	numWps := 100
	jobs := make(chan []string, numWps)
	res := make(chan *api.Order)

	reprocess := make([]*api.Order, 0)

	var wg sync.WaitGroup
	worker := func(jobs <-chan []string, results chan<- *api.Order) {
		for {
			select {
			case job, ok := <-jobs: // you must check for readable state of the channel.
				if !ok {
					return
				}
				result, err := parseStruct(job)
				if err != nil {
					fmt.Println(err)
					return
				}
				errPost := callPersistenceApi(result)
				if errPost != nil {
					if errPost.Error() == "couldn't create resource" {
						// If couldn't create resource save the object to retry later
						results <- result
					} else {
						fmt.Println(errPost)
					}
				}
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
		reprocess = append(reprocess, r)
	}

	fmt.Println(len(reprocess))
}

func callPersistenceApi(order *api.Order) error {
	json, err := json.Marshal(order)
	if err != nil {
		return err
	}
	response, err := http.Post("https://localhost:9101/orders", "application/json", bytes.NewBuffer(json))

	if err != nil {
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
	phoneNumber := strings.TrimSpace(data[2])
	parcelWeight, err := strconv.ParseFloat(strings.TrimSpace(data[3]), 32)
	if err != nil {
		fmt.Println(err)
	}

	year, month, day := time.Now().Date()
	date := strconv.Itoa(year) + "-" + month.String() + "-" + strconv.Itoa(day)

	country, err := getCountry(phoneNumber)

	if err != nil {
		return nil, err
	}

	return api.CreateOrder(int(id), email, phoneNumber, float32(parcelWeight), date, country), nil
}

func getCountry(phoneNumber string) (string, error) {
	return "Brazil", nil
}
