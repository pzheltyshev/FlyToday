package flight

import (
	"context"
	"log/slog"
	"time"
)

type FlightService struct {
	repo     Repository
	provider FlightProvider
	logger   *slog.Logger
}

func NewFlightService(repo Repository, provider FlightProvider, logger *slog.Logger) *FlightService {
	return &FlightService{
		repo:     repo,
		provider: provider,
		logger:   logger,
	}
}

func (f *FlightService) SearchFlight(ctx context.Context, departureAirport string, arrivalAirport string, date time.Time) Flight {

	flightSegments, err := f.provider.SearchFlight(ctx, departureAirport, arrivalAirport, date)

	if err != nil {
		return Flight{}
	}

	//TODO add test on correct flight data (like airport code)

	err = f.repo.SaveRequest(ctx, flightSegments)
	if err != nil {
		f.logger.Error("Error while saving request")
	}

	flight := f.ConvSegmentsRawToFlight(ctx, flightSegments)

	return flight

}

func (f *FlightService) GetAirportRefByIAITCode(ctx context.Context, code string) AirportRef {

	airport, err := f.repo.GetAirportRefByIAITCode(ctx, code)

	if err != nil {
		return AirportRef{}
	}

	return airport

}

func (f *FlightService) GetAirlineRefFromICAO(ctx context.Context, code string) AirlineRef {

	airline, err := f.repo.GetAirportRefByICAOCode(ctx, code)

	if err != nil {
		return AirlineRef{}
	}

	return airline

}

func (f *FlightService) ConvSegmentsRawToFlight(ctx context.Context, SegmentsRaw *FlightSegmentsRaw) Flight {

	flight := Flight{}

	flight.Price = SegmentsRaw.Price
	flight.Currency = SegmentsRaw.Currency

	for _, data := range SegmentsRaw.Segments {

		origin := f.GetAirportRefByIAITCode(ctx, data.Origin)
		destination := f.GetAirportRefByIAITCode(ctx, data.Destination)
		airlineRef := f.GetAirlineRefFromICAO(ctx, data.AirlineCode)

		flight.Segments = append(flight.Segments, FlightSegment{
			DepartureAirport: origin,
			ArrivalAirport:   destination,
			DepartureTime:    data.DateFrom,
			ArrivalTime:      data.DateTo,
			Airline:          airlineRef,
		})

	}
}
