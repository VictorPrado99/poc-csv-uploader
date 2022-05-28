package api

type Order struct {
	Id           int     `json:"id"`
	Email        string  `json:"email"`
	PhoneNumber  string  `json:"phone_number"`
	ParcelWeight float32 `json:"parcel_weight"`
	Date         string  `json:"date"`
	Country      string  `json:"country"`
}

func CreateOrder(id int, email string, phoneNumber string, parcelWeight float32, date string, country string) *Order {
	return &Order{
		id,
		email,
		phoneNumber,
		parcelWeight,
		date,
		country,
	}
}
