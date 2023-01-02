CREATE TABLE IF NOT EXISTS entries (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    rating INTEGER NOT NULL,
    testimonial TEXT NOT NULL,
    last_name TEXT NOT NULL,
    first_name TEXT NOT NULL,
    occupation TEXT NOT NULL,
    github TEXT NOT NULL
);

ALTER TABLE entries ADD CONSTRAINT entry_rating_check CHECK (rating BETWEEN 1 AND 5);