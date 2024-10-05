CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS shortened_urls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    original_url TEXT NOT NULL,
    hash VARCHAR(7) NOT NULL,
    number_of_accesses INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX shortened_urls_hash_index ON shortened_urls (hash);
CREATE UNIQUE INDEX shotened_urls_original_url_index ON shortened_urls (original_url);