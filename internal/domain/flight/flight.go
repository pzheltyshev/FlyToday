package flight

import "time"

type flightSegment struct {
	id               int
	departureAirport Airport
	arrivalAirport   Airport
	departureTime    time.Time
	arrivalTime      time.Time
	airline          airline
	airlineFlight    airlineFlight
}

type Flight struct {
	id               int
	departureAirport Airport
	arrivalAirport   Airport
	price            float32
	currency         string
	segments         []flightSegment
}

type Airport struct {
	IATAcode string
	name     string
	city     string
}

type airline struct {
	code    string
	name    string
	flights []airlineFlight
}

type airlineFlight struct {
	number string
}
