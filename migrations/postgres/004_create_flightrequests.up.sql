CREATE TABLE flight_requests(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    request_date TIMESTAMPTZ,

    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)

)