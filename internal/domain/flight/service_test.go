package flight

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"reflect"
	"strings"
	"testing"
	"time"
)

type fakeFlightProvider struct {
	result           *FlightSegmentsRaw
	err              error
	called           bool
	ctx              context.Context
	departureAirport string
	arrivalAirport   string
	date             time.Time
}

func (p *fakeFlightProvider) SearchFlight(ctx context.Context, departureAirport string, arrivalAirport string, date time.Time) (*FlightSegmentsRaw, error) {
	p.called = true
	p.ctx = ctx
	p.departureAirport = departureAirport
	p.arrivalAirport = arrivalAirport
	p.date = date
	return p.result, p.err
}

type fakeFlightRepository struct {
	saveErr    error
	saved      *FlightSegmentsRaw
	airport    AirportRef
	airportErr error
	airline    AirlineRef
	airlineErr error
}

func (r *fakeFlightRepository) SaveRequest(ctx context.Context, request *FlightSegmentsRaw) error {
	r.saved = request
	return r.saveErr
}

func (r *fakeFlightRepository) GetAirportRefByIATACode(ctx context.Context, code string) (AirportRef, error) {
	return r.airport, r.airportErr
}

func (r *fakeFlightRepository) GetAirlineRefFromICAO(ctx context.Context, code string) (AirlineRef, error) {
	return r.airline, r.airlineErr
}

func newTestLogger(writer ...io.Writer) *slog.Logger {
	var output io.Writer = io.Discard
	if len(writer) > 0 {
		output = writer[0]
	}
	return slog.New(slog.NewTextHandler(output, nil))
}

func TestFlightServiceSearchFlightSuccess(t *testing.T) {
	ctx := context.Background()
	date := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	raw := &FlightSegmentsRaw{
		RequestDate: date,
		Price:       125.5,
		Currency:    "USD",
		Segments: []FlightSegmentRaw{
			{
				Origin:      "LED",
				Destination: "PVG",
				DateFrom:    time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC),
				DateTo:      time.Date(2026, 10, 1, 13, 0, 0, 0, time.UTC),
				AirlineCode: "MU",
			},
			{
				Origin:      "PVG",
				Destination: "HKG",
				DateFrom:    time.Date(2026, 10, 1, 14, 0, 0, 0, time.UTC),
				DateTo:      time.Date(2026, 10, 1, 17, 0, 0, 0, time.UTC),
				AirlineCode: "CX",
			},
		},
	}
	repository := &fakeFlightRepository{
		airport: AirportRef{IATACode: "LED", Name: "Pulkovo"},
		airline: AirlineRef{Code: "MU", Name: "China Eastern"},
	}
	provider := &fakeFlightProvider{result: raw}
	service := NewFlightService(repository, provider, newTestLogger())

	result := service.SearchFlight(ctx, "LED", "HKG", date)

	if !provider.called {
		t.Fatal("provider was not called")
	}
	if provider.departureAirport != "LED" || provider.arrivalAirport != "HKG" || !provider.date.Equal(date) {
		t.Errorf("provider received wrong arguments: %#v", provider)
	}

	expected := Flight{
		Price:    125.5,
		Currency: "USD",
		Segments: []FlightSegment{
			{
				DepartureAirport: AirportRef{IATACode: "LED", Name: "Pulkovo"},
				ArrivalAirport:   AirportRef{IATACode: "LED", Name: "Pulkovo"},
				DepartureTime:    raw.Segments[0].DateFrom,
				ArrivalTime:      raw.Segments[0].DateTo,
				Airline:          AirlineRef{Code: "MU", Name: "China Eastern"},
			},
			{
				DepartureAirport: AirportRef{IATACode: "LED", Name: "Pulkovo"},
				ArrivalAirport:   AirportRef{IATACode: "LED", Name: "Pulkovo"},
				DepartureTime:    raw.Segments[1].DateFrom,
				ArrivalTime:      raw.Segments[1].DateTo,
				Airline:          AirlineRef{Code: "MU", Name: "China Eastern"},
			},
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("result = %#v, want %#v", result, expected)
	}
}

func TestFlightServiceSearchFlightProviderError(t *testing.T) {
	ctx := context.Background()
	date := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	repository := &fakeFlightRepository{}
	provider := &fakeFlightProvider{err: errors.New("provider failed")}
	service := NewFlightService(repository, provider, newTestLogger())

	result := service.SearchFlight(ctx, "LED", "PVG", date)

	if !provider.called {
		t.Fatal("provider was not called")
	}

	if !reflect.DeepEqual(result, Flight{}) {
		t.Errorf("result = %#v, want zero Flight", result)
	}
}

func TestFlightServiceSearchFlightSaveRequestError(t *testing.T) {
	ctx := context.Background()
	date := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	raw := &FlightSegmentsRaw{
		Price:    125.5,
		Currency: "USD",
		Segments: []FlightSegmentRaw{
			{
				Origin:      "LED",
				Destination: "PVG",
				DateFrom:    time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC),
				DateTo:      time.Date(2026, 10, 1, 13, 0, 0, 0, time.UTC),
				AirlineCode: "MU",
			},
		},
	}
	var logs bytes.Buffer
	repository := &fakeFlightRepository{
		saveErr: errors.New("save failed"),
		airport: AirportRef{IATACode: "LED", Name: "Pulkovo"},
		airline: AirlineRef{Code: "MU", Name: "China Eastern"},
	}
	provider := &fakeFlightProvider{result: raw}
	service := NewFlightService(repository, provider, newTestLogger(&logs))

	result := service.SearchFlight(ctx, "LED", "PVG", date)

	if !strings.Contains(logs.String(), "Error while saving request") {
		t.Errorf("logger output = %q, want save error message", logs.String())
	}
	if result.Price != raw.Price || result.Currency != raw.Currency || len(result.Segments) != 1 {
		t.Errorf("result = %#v, want converted flight", result)
	}
}

func TestFlightServiceConvSegmentsRawToFlightWithoutSegments(t *testing.T) {
	ctx := context.Background()
	repository := &fakeFlightRepository{}
	provider := &fakeFlightProvider{}
	service := NewFlightService(repository, provider, newTestLogger())
	raw := &FlightSegmentsRaw{Price: 10, Currency: "RUB"}

	result := service.ConvSegmentsRawToFlight(ctx, raw)

	if result.Price != 10 || result.Currency != "RUB" {
		t.Errorf("result = %#v, want price and currency", result)
	}
	if len(result.Segments) != 0 {
		t.Errorf("segments = %d, want 0", len(result.Segments))
	}
}

func TestFlightServiceGetAirportRefByIATACode(t *testing.T) {
	expected := AirportRef{IATACode: "LED", Name: "Pulkovo"}
	repository := &fakeFlightRepository{airport: expected}
	service := NewFlightService(repository, &fakeFlightProvider{}, newTestLogger())

	result := service.GetAirportRefByIATACode(context.Background(), "LED")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("result = %#v, want %#v", result, expected)
	}
}

func TestFlightServiceGetAirportRefByIATACodeError(t *testing.T) {
	repository := &fakeFlightRepository{airportErr: errors.New("airport lookup failed")}
	service := NewFlightService(repository, &fakeFlightProvider{}, newTestLogger())

	result := service.GetAirportRefByIATACode(context.Background(), "LED")

	if !reflect.DeepEqual(result, AirportRef{}) {
		t.Errorf("result = %#v, want zero AirportRef", result)
	}
}

func TestFlightServiceGetAirlineRefFromICAO(t *testing.T) {
	expected := AirlineRef{Code: "MU", Name: "China Eastern"}
	repository := &fakeFlightRepository{airline: expected}
	service := NewFlightService(repository, &fakeFlightProvider{}, newTestLogger())

	result := service.GetAirlineRefFromICAO(context.Background(), "MU")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("result = %#v, want %#v", result, expected)
	}
}

func TestFlightServiceGetAirlineRefFromICAOError(t *testing.T) {
	repository := &fakeFlightRepository{airlineErr: errors.New("airline lookup failed")}
	service := NewFlightService(repository, &fakeFlightProvider{}, newTestLogger())

	result := service.GetAirlineRefFromICAO(context.Background(), "MU")

	if !reflect.DeepEqual(result, AirlineRef{}) {
		t.Errorf("result = %#v, want zero AirlineRef", result)
	}
}
