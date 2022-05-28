package logic

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/VictorPrado99/poc-csv-uploader/pkg/api"
)

func ProcessCsv(tempFile *os.File) {
	tempFile, _ = os.Open(tempFile.Name())

	fcsv := csv.NewReader(tempFile)
	rs := make([]*api.Order, 0)
	numWps := 100
	jobs := make(chan []string, numWps)
	res := make(chan *api.Order)

	var wg sync.WaitGroup
	worker := func(jobs <-chan []string, results chan<- *api.Order) {
		for {
			select {
			case job, ok := <-jobs: // you must check for readable state of the channel.
				if !ok {
					return
				}
				results <- parseStruct(job)
			}
		}
	}

	// init workers
	for w := 0; w < numWps; w++ {
		wg.Add(1)
		go func() {
			// this line will exec when chan `res` processed output at line 107 (func worker: line 71)
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
		rs = append(rs, r)
	}

	fmt.Println("Count Concu ", len(rs))

}

func callPersistenceApi(order *api.Order) {

}

func parseStruct(data []string) *api.Order {
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
		fmt.Println(err)
	}

	return api.CreateOrder(int(id), email, phoneNumber, float32(parcelWeight), date, country)
}

func getCountry(phoneNumber string) (string, error) {
	return "Brazil", nil
}
