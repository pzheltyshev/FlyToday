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
	Flights     []flightProvider
}

type flightProvider struct {
	Id       int
	Provider string
	Price    float32
	Currency string
	Segments []segment
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

func (f *flightRequestRow) toDomain() []flight.Flight {

	result := make([]flight.Flight, 0, len(f.Flights))

	for i, data := range f.Flights {

		result = append(result, flight.Flight{
			Currency: data.Currency,
			Price:    data.Price,
			Provider: data.Provider,
			Segments: make([]flight.FlightSegment, 0, len(data.Segments)),
		})

		for _, s := range data.Segments {
			result[i].Segments = append(result[i].Segments, flight.FlightSegment{

				DepartureAirport: flight.AirportRef{
					IATACode: s.OriginIATAcode,
				},
				ArrivalAirport: flight.AirportRef{
					IATACode: s.DestinationIATAcode,
				},
				DepartureTime: s.DateFrom,
				ArrivalTime:   s.DateTo,
				Airline: flight.AirlineRef{
					Code: s.AirlineCode,
				},
			})
		}
	}

	return result
}

func (f *flightRequestRow) fromRequest(date time.Time, request []flight.FlightSegmentsRaw) {

	f.RequestDate = date

	for i, data := range request {

		f.Flights = append(f.Flights, flightProvider{

			Price:    data.Price,
			Currency: data.Currency,
			Provider: data.Provider,
			Segments: make([]segment, 0, len(data.Segments)),
		})

		for _, s := range data.Segments {
			f.Flights[i].Segments = append(f.Flights[i].Segments, segment{
				OriginIATAcode:      s.Origin,
				DestinationIATAcode: s.Destination,
				DateFrom:            s.DateFrom,
				DateTo:              s.DateTo,
				AirlineCode:         s.AirlineCode,
				FlightCode:          s.FlightCode,
			})
		}

	}

}

func (f *FlightRepository) SaveRequest(ctx context.Context, requestDate time.Time, flightSegments []flight.FlightSegmentsRaw) error {

	flightRequest := NewFlightRequestRow()
	flightRequest.fromRequest(requestDate, flightSegments)

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
			request_date
		)
		VALUES($1, $2)
		RETURNING id
	`,
		flightRequest.UserId,
		flightRequest.RequestDate,
	).Scan(&requestId)

	if err != nil {
		return err
	}

	flightRequest.Id = requestId

	for _, data := range flightRequest.Flights {

		err = tx.QueryRowContext(ctx, `
			INSERT INTO flights(
				price,
				currency,
				provider
			)
			VALUES($1, $2, $3)
			RETURNING id
		`,
			data.Currency,
			data.Price,
			data.Provider,
		).Scan(&requestId)

		if err != nil {
			return err
		}

		data.Id = requestId

		for i, s := range data.Segments {
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
				s.OriginIATAcode,
				s.DestinationIATAcode,
				s.DateFrom,
				s.DateTo,
				s.AirlineCode,
				i,
			)

			if err != nil {
				return err
			}
		}

	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (f *FlightRepository) GetAirportRefByIATACode(ctx context.Context, code string) (flight.AirportRef, error) {

	airportRef := flight.AirportRef{}

	err := f.db.QueryRowContext(ctx, `
		SELECT iata_code, icao_code
		FROM airports
		WHERE iata_code = $1
	`,
		code,
	).Scan(
		&airportRef.IATACode,
		&airportRef.Name,
	)

	if err != nil {
		return flight.AirportRef{}, err
	}

	return airportRef, nil
}

func (f *FlightRepository) GetAirlineRefFromICAO(ctx context.Context, code string) (flight.AirlineRef, error) {

	airlineRef := flight.AirlineRef{}

	err := f.db.QueryRowContext(ctx, `
		SELECT icao_code, name
		FROM airlines
		WHERE icao_code = $1
	`,
		code,
	).Scan(
		&airlineRef.Code,
		&airlineRef.Name,
	)

	if err != nil {
		return flight.AirlineRef{}, err
	}

	return airlineRef, nil

}
