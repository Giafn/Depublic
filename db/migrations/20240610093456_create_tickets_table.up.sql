BEGIN;

CREATE TABLE tickets (
    id UUID PRIMARY KEY,
    id_transaction VARCHAR(255) NOT NULL,
    id_event VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    booking_num VARCHAR(255) NOT NULL,
    is_used BOOLEAN NOT NULL
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL
);

COMMIT;