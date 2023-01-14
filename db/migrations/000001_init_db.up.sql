CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT
);

CREATE TABLE IF NOT EXISTS cred_data (
    id            SERIAL PRIMARY KEY,
    meta          TEXT UNIQUE NOT NULL,
    email         BYTEA,
    password_hash BYTEA
);

CREATE TABLE IF NOT EXISTS bin_data (
    id           SERIAL PRIMARY KEY,
    meta         TEXT UNIQUE NOT NULL,
    bytes        BYTEA
);

CREATE TABLE IF NOT EXISTS text_data (
    id           SERIAL PRIMARY KEY,
    meta         TEXT UNIQUE NOT NULL,
    text         BYTEA
);

CREATE TABLE IF NOT EXISTS card_data (
    id            SERIAL PRIMARY KEY,
    meta          TEXT UNIQUE NOT NULL,
    num           BYTEA,
    period_dt     BYTEA,
    cvv           BYTEA,
    full_name     BYTEA
);
