BEGIN;

CREATE TABLE tickets (
    id UUID PRIMARY KEY,
    transaction_id UUID NOT NULL,
    event_id UUID NOT NULL,
    pricing_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    booking_num VARCHAR(255) NOT NULL,
    is_used BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL
);

COMMIT;