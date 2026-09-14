package mock

import (
	"context"
	"time"

	"github.com/pzheltyshev/FlyToday/internal/domain/flight"
)

const provider = "mock_provider"

type FlightProvider struct {
}

func NewFlightProvider() *FlightProvider {

	return &FlightProvider{}
}

func (f *FlightProvider) SearchFlight(ctx context.Context, departureAirport string, arrivalAirport string, date time.Time) ([]flight.FlightSegmentsRaw, error) {

	var segments []flight.FlightSegmentsRaw

	segments = append(segments, flight.FlightSegmentsRaw{
		RequestDate: time.Date(2026, 8, 23, 20, 10, 0, 0, time.UTC),
		Provider:    provider,
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
	})

	return segments, nil

}
