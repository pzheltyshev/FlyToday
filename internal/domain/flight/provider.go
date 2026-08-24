package flight

import (
	"context"
	"time"
)

type FlightProvider interface {
	SearchFlight(ctx context.Context, departureAirport string, arrivalAirport string, date time.Time) (*FlightSegmentsRaw, error)
}

type FlightSegmentRaw struct {
	Origin      string
	Destination string
	DateFrom    time.Time
	DateTo      time.Time
	AirlineCode string
}

type FlightSegmentsRaw struct {
	ResponseDate time.Time
	Segments     []FlightSegmentRaw
}
