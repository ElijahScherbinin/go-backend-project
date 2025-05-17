DROP DATABASE IF EXISTS user_service;
CREATE DATABASE user_service;

DROP TABLE IF EXISTS users;

\c user_service;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(32) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE
);

CREATE INDEX idx_username ON users (username);
CREATE INDEX idx_email ON users (email);