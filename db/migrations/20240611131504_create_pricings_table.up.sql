BEGIN;
CREATE TABLE IF NOT EXISTS pricings (
    pricing_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    quota INTEGER NOT NULL,
    remaining INTEGER NOT NULL,
    fee INTEGER NOT NULL
);
COMMIT;
