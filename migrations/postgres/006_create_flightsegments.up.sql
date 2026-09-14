CREATE TABLE flight_segments(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    flight_id BIGINT NOT NULL,
    departure_airport_id VARCHAR(3) NOT NULL,
    arrival_airport_id VARCHAR(3) NOT NULL,
    departure_at TIMESTAMPTZ NOT NULL,
    arrival_at TIMESTAMPTZ NOT NULL,
    airline_code VARCHAR(3) NOT NULL,
    segment_order INT NOT NULL,

    CONSTRAINT fk_flights
        FOREIGN KEY (flight_id)
        REFERENCES flights(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_departure_airport
        FOREIGN KEY (departure_airport_id)
        REFERENCES airports
        ON DELETE CASCADE,

    CONSTRAINT fk_arrival_airport
        FOREIGN KEY (arrival_airport_id)
        REFERENCES airports
        ON DELETE CASCADE,

    CONSTRAINT fk_airline
        FOREIGN KEY (airline_code)
        REFERENCES airlines
        ON DELETE CASCADE,

    CONSTRAINT unique_flight_segment_order
        UNIQUE (flight_request_id, segment_order)
)