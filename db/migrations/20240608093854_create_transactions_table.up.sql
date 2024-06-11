BEGIN;

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
    user_id UUID NOT NULL,
    ticket_quantity INTEGER NOT NULL,
    total_amount INTEGER NOT NULL,
    is_paid BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT fk_event FOREIGN KEY (event_id) REFERENCES events(event_id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(user_id)
);

ALTER TABLE transactions ADD COLUMN auditable JSONB;

COMMIT;
