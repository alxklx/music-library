CREATE SCHEMA IF NOT EXISTS "music-library";

CREATE TABLE IF NOT EXISTS "music-library".songs (
    id SERIAL PRIMARY KEY,
    "group" VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL,
    release_date VARCHAR(50),
    text TEXT,
    link VARCHAR(255)
);