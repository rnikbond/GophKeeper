CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    email         CHARACTER VARYING(50),
    password_hash CHARACTER VARYING(64)
);
