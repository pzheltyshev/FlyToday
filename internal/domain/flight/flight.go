package flight

import "time"

type FlightSegment struct {
	DepartureAirport AirportRef
	ArrivalAirport   AirportRef
	DepartureTime    time.Time
	ArrivalTime      time.Time
	Airline          AirlineRef
}

type Flight struct {
	Id       int
	Price    float32
	Currency string
	Segments []FlightSegment
}

type AirportRef struct {
	IATACode string
	Name     string
}

type AirlineRef struct {
	Code       string
	Name       string
	FlightCode string
}
