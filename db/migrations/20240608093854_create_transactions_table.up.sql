BEGIN;

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
    user_id UUID NOT NULL,
    ticket_quantity INTEGER NOT NULL,
    total_amount INTEGER NOT NULL,
    is_paid BOOLEAN DEFAULT FALSE,
    payment_url TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT fk_event_id
        FOREIGN KEY (event_id)
        REFERENCES events (id)
        ON DELETE CASCADE
);

COMMIT;
