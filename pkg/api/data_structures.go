package api

type DataTuple struct {
	Id           int
	Email        string
	PhoneNumber  string
	ParcelWeight float32
	Date         string
	Country      string
}

func CreateTuple(id int, email string, phoneNumber string, parcelWeight float32, date string, country string) *DataTuple {
	return &DataTuple{
		id,
		email,
		phoneNumber,
		parcelWeight,
		date,
		country,
	}
}
