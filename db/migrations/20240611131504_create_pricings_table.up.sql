BEGIN;
CREATE TABLE IF NOT EXISTS pricings (
    pricing_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    fee INTEGER NOT NULL,
    CONSTRAINT fk_event FOREIGN KEY (event_id) REFERENCES events(id)
);
COMMIT;
