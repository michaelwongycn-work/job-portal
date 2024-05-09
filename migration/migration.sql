CREATE DATABASE jobPortal;
USE jobPortal;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role INTEGER NOT NULL
);

CREATE TABLE user_tokens (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    access_token VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    expiration_time BIGINT NOT NULL,
    PRIMARY KEY (user_id)
);

CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    employer_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    requirement TEXT NOT NULL,
    create_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE applications (
    id SERIAL PRIMARY KEY,
    job_id INTEGER REFERENCES jobs(id) ON DELETE CASCADE,
    talent_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    application_status INTEGER NOT NULL DEFAULT 1,
    apply_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    UNIQUE (job_id, talent_id)
);
