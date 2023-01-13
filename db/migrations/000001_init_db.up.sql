CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    email         CHARACTER VARYING(50) UNIQUE NOT NULL,
    password_hash CHARACTER VARYING(64)
);

CREATE TABLE IF NOT EXISTS cred_data (
    id            SERIAL PRIMARY KEY,
    meta          CHARACTER VARYING(50) UNIQUE NOT NULL,
    email         CHARACTER VARYING(50),
    password_hash CHARACTER VARYING(64)
);

CREATE TABLE IF NOT EXISTS bin_data (
    id           SERIAL PRIMARY KEY,
    meta         CHARACTER VARYING(50) UNIQUE NOT NULL,
    bytes        bytea
);

CREATE TABLE IF NOT EXISTS text_data (
    id           SERIAL PRIMARY KEY,
    meta         CHARACTER VARYING(50) UNIQUE NOT NULL,
    text         TEXT
);

CREATE TABLE IF NOT EXISTS card_data (
    id            SERIAL PRIMARY KEY,
    meta          CHARACTER VARYING(50) UNIQUE NOT NULL,
    num           CHARACTER VARYING(16),
    period_dt     CHARACTER VARYING(7),
    cvv           CHARACTER VARYING(3),
    full_name     TEXT
);
