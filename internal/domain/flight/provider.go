package flight

import (
	"context"
	"time"
)

type FlightProvider interface {
	SearchFlight(ctx context.Context, departureAirport string, arrivalAirport string, date time.Time) (*FlightSegmentsRaw, error)
}
