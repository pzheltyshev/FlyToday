package flight

import (
	"context"
)

type Repository interface {
	SaveRequest(ctx context.Context, flightSegments FlightSegmentsRaw) error
}
