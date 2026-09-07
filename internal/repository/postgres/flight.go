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
	UserId      int
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

func (f *flightRequestRow) fromRequest(request *flight.FlightSegmentsRaw) {

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

func (f *FlightRepository) SaveRequest(ctx context.Context, flightSegments *flight.FlightSegmentsRaw) error {

	flightRequest := NewFlightRequestRow()
	flightRequest.fromRequest(flightSegments)

	tx, err := f.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var requestId int

	err = tx.QueryRowContext(ctx, `
		INSERT INTO flight_requests(
			user_id,
			price,
			currency
		)
		VALUES($1, $2, $3)
		RETURNING id
	`,
		flightRequest.UserId,
		flightRequest.Price,
		flightRequest.Currency,
	).Scan(&requestId)

	if err != nil {
		return err
	}

	flightRequest.Id = requestId

	for i, data := range flightRequest.Segments {
		_, err = tx.QueryContext(ctx, `
			INSERT INTO flight_segments(
				flight_request_id,
				departure_airport_id,
				arrival_airport_id,
				departure_at,
				arrival_at,
				airline_code,
				segment_order
			)
		`,
			flightRequest.Id,
			data.OriginIATAcode,
			data.DestinationIATAcode,
			data.DateFrom,
			data.DateTo,
			data.AirlineCode,
			i,
		)

		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (f *FlightRepository) GetAirportRefByIAITCode(ctx context.Context, code string) flight.AirportRef {

	airportRef := flight.AirportRef{}

	err := f.db.QueryRowContext(ctx, `
		SELECT iata_code, icao_code
		FROM airports
		WHERE iata_code = $1
	`,
		code,
	).Scan(
		airportRef.IATACode,
		airportRef.Name,
	)

	if err != nil {
		return flight.AirportRef{}
	}

	return airportRef
}

func (f *FlightRepository) GetAirlineRefFromICAO(ctx context.Context, code string) flight.AirlineRef {

	airlineRef := flight.AirlineRef{}

	err := f.db.QueryRowContext(ctx, `
		SELECT icao_code, name
		FROM airlines
		WHERE icao_code = $1
	`,
		code,
	).Scan(
		airlineRef.Code,
		airlineRef.Name,
	)

	if err != nil {
		return flight.AirlineRef{}
	}

	return airlineRef

}
