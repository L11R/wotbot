CREATE EXTENSION IF NOT EXISTS "moddatetime";

CREATE TABLE IF NOT EXISTS users
(
    id           BIGSERIAL PRIMARY KEY,
    telegram_id  INTEGER NOT NULL UNIQUE,
    nickname     TEXT,
    wargaming_id INTEGER,
    created_at   TIMESTAMP DEFAULT now(),
    updated_at   TIMESTAMP
);

CREATE TRIGGER update_updated_at
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE moddatetime(updated_at);

CREATE TABLE IF NOT EXISTS stats
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    name       TEXT    NOT NULL,
    value      TEXT    NOT NULL,
    html_id    TEXT    NOT NULL,
    trend_img  BYTEA   NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);