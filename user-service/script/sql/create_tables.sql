DROP DATABASE IF EXISTS user_service;
CREATE DATABASE user_service;

DROP TABLE IF EXISTS users;

\c user_service;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(32) NOT NULL UNIQUE,
    password_hash VARCHAR(64) NOT NULL
);

CREATE INDEX idx_username ON users (username);