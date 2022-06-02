package csvuploaderlogic

import (
	"os"
	"strings"
	"testing"

	"github.com/VictorPrado99/poc-csv-persistence/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CsvProcessingTestSuite struct {
	suite.Suite
}

func TestCsvProcessingSuite(t *testing.T) {
	suite.Run(t, new(CsvProcessingTestSuite))
}

func (suite *CsvProcessingTestSuite) TestParseStructSuccess() {
	data := []string{"1", "test@test.com", "237 209993809", "24.45"}

	got, err := parseStruct(data)

	if err != nil {
		suite.T().Fatalf("Couldn't correctly parse, %s \n", strings.Join(data, ", "))
	}

	want := &api.Order{
		Id:           1,
		Email:        "test@test.com",
		PhoneNumber:  "(237)209993809",
		ParcelWeight: 24.45,
		Date:         getDate(),
		Country:      "Cameroon",
	}

	assert.Equal(suite.T(), want, got)
}

func (suite *CsvProcessingTestSuite) TestParseStructFail() {
	data1 := []string{"1", "test@test.com", "237 209993809", "24.dx45"}

	_, err1 := parseStruct(data1)

	if err1 == nil {
		suite.T().Fatalf("Failed to throw error when couldn't parse to float with data %s\n", data1[3])
	}

	data2 := []string{"1", "test@test.com", "55 209993809", "24.45"}

	_, err2 := parseStruct(data2)

	if err2 == nil {
		suite.T().Fatalf("Failed to throw error when with unrecognized country number, %s \n", data2[2])
	}

}

func (suite *CsvProcessingTestSuite) TestDoCsvProcessingSuccess() {
	file, err := os.Open("../../test_files/success_test.csv")

	if err != nil {
		suite.T().Fatalf("Couldn't open csv %v", err)
	}

	got := DoCsvProcess(file)
	dWant := []*api.Order{}

	assert.NotEqual(suite.T(), dWant, got)
}

func (suite *CsvProcessingTestSuite) TestDoCsvProcessingFail() {
	file, err := os.Open("../../test_files/fail_test.csv")

	if err != nil {
		suite.T().Fatalf("Couldn't open csv %v", err)
	}

	got := DoCsvProcess(file)
	want := []*api.Order{}

	assert.Equal(suite.T(), want, got)
}
