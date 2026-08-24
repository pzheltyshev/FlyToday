package postgres

import (
	"time"

	"github.com/pzheltyshev/FlyToday/internal/domain/flight"
)

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
	OriginName          string
	DestinationIATAcode string
	DestinationName     string
	DateFrom            time.Time
	DateTo              time.Time
	AirlineCode         string
	AirlineName         string
	FlightCode          string
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
