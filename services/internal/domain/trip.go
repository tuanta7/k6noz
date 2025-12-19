package domain

type Trip struct {
	ID          string      `json:"id" `
	Location    []*Location `json:"location"`
	DriverID    string      `json:"driverId"`
	PassengerID string      `json:"passengerId"`
}
