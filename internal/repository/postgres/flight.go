package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/pzheltyshev/FlyToday/internal/domain/flight"
)

type FlightRepository struct {
	db *sql.DB
}

func NewFlightRepository(db *sql.DB) *FlightRepository {
	return &FlightRepository{
		db: db,
	}
}

type flightRequestRow struct {
	Id          int
	IdUser      int
	UserName    string
	RequestDate time.Time
	Price       float32
	Currency    string
	Segments    []segment
}

type segment struct {
	OriginIATAcode      string
	DestinationIATAcode string
	DateFrom            time.Time
	DateTo              time.Time
	AirlineCode         string
	FlightCode          string
}

func NewFlightRequestRow() *flightRequestRow {
	return &flightRequestRow{}
}

type airportRow struct {
}

type airlineRow struct {
}

func (f *flightRequestRow) toDomain() flight.Flight {

	flightData := flight.Flight{}

	flightData.Id = f.Id
	flightData.Currency = f.Currency
	flightData.Price = f.Price

	for _, data := range f.Segments {

		originRef := flight.AirportRef{
			IATACode: data.OriginIATAcode,
			Name:     data.OriginName,
		}

		destinationRef := flight.AirportRef{
			IATACode: data.DestinationIATAcode,
			Name:     data.DestinationName,
		}

		airline := flight.AirlineRef{
			Code:       data.AirlineCode,
			Name:       data.AirlineName,
			FlightCode: data.FlightCode,
		}

		flightData.Segments = append(flightData.Segments, flight.FlightSegment{
			DepartureAirport: originRef,
			ArrivalAirport:   destinationRef,
			DepartureTime:    data.DateFrom,
			ArrivalTime:      data.DateTo,
			Airline:          airline,
		})
	}

	return flightData
}

func (f *flightRequestRow) fromRequest(request flight.FlightSegmentsRaw) {

	f.RequestDate = request.RequestDate
	f.Price = request.Price
	f.Currency = request.Currency

	for _, data := range request.Segments {
		f.Segments = append(f.Segments, segment{
			OriginIATAcode:      data.Origin,
			DestinationIATAcode: data.Destination,
			DateFrom:            data.DateFrom,
			DateTo:              data.DateTo,
			AirlineCode:         data.AirlineCode,
			FlightCode:          data.FlightCode,
		})
	}

}

func (f *FlightRepository) SaveRequest(ctx context.Context, flightSegments flight.FlightSegmentsRaw) error {

	flightRequest := NewFlightRequestRow()
	flightRequest.fromRequest(flightSegments)

	f.db.Query("")

	return nil
}
