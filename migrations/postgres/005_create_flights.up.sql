CREATE TABLE flights(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    price REAL,
    currency VARCHAR(3),
    provider VARCHAR(100),
    flight_request_id BIGINT NOT NULL

    CONSTRAINT fk_flight_request
        FOREIGN KEY (flight_request_id)
        REFERENCES flight_requests(id)

)