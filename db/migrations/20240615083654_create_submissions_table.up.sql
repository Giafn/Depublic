BEGIN;

CREATE TABLE IF NOT EXISTS submissions(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL,
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    status varchar(50),
    type varchar(50) NOT NULL
);

COMMIT;