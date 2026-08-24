package mock

import (
	"context"
	"time"

	"github.com/pzheltyshev/FlyToday/internal/domain/flight"
)

type FlightProvider struct {
}

func NewFlightProvider() *FlightProvider {

	return &FlightProvider{}
}

func (f *FlightProvider) SearchFlight(ctx context.Context, departureAirport string, arrivalAirport string, date time.Time) (*flight.FlightSegmentsRaw, error) {

	segments := flight.FlightSegmentsRaw{
		ResponseDate: time.Date(2026, 8, 23, 20, 10, 0, 0, time.UTC),
		Segments: []flight.FlightSegmentRaw{
			{
				Origin:      "LED",
				Destination: "PVG",
				DateFrom:    time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC),
				DateTo:      time.Date(2026, 10, 1, 13, 0, 0, 0, time.UTC),
				AirlineCode: "MU",
			},
			{
				Origin:      "LED",
				Destination: "PVG",
				DateFrom:    time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC),
				DateTo:      time.Date(2026, 10, 1, 13, 0, 0, 0, time.UTC),
				AirlineCode: "MU",
			},
		},
	}

	return &segments, nil

}
