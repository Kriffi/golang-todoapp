CREATE SCHEMA todoapp;

CREATE TABLE todoapp.users (
    id           SERIAL                PRIMARY KEY,
    version      BIGINT       NOT NULL DEFAULT 1,
    full_name    VARCHAR(100) NOT NULL CHECK (char_jenght(full_name) BETWEEN 3 AND 100),
    phone_number VARCHAR(15)
);