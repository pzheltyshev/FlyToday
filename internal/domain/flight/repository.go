package flight

import (
	"context"
	"time"
)

type Repository interface {
	SaveRequest(ctx context.Context, requestDate time.Time, flightSegments []FlightSegmentsRaw) error
	GetAirportRefByIATACode(ctx context.Context, code string) (AirportRef, error)
	GetAirlineRefFromICAO(ctx context.Context, code string) (AirlineRef, error)
}
