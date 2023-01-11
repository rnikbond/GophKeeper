CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    email         CHARACTER VARYING(50) UNIQUE NOT NULL,
    password_hash CHARACTER VARYING(64)
);

CREATE TABLE IF NOT EXISTS text_data (
    id           SERIAL PRIMARY KEY,
    meta         CHARACTER VARYING(50) UNIQUE NOT NULL,
    text         TEXT
);