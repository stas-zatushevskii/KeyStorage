package db

import "database/sql"

var createTable = `
DROP TABLE IF EXISTS user_tokens CASCADE;
DROP TABLE IF EXISTS bank_data CASCADE;
DROP TABLE IF EXISTS file_data CASCADE;
DROP TABLE IF EXISTS text_data CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS account_data CASCADE;

CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE text_data (
    id      BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,

    title   VARCHAR(256),
    text    TEXT,

    CONSTRAINT fk_text_data_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

CREATE TABLE file_data (
    id        BIGSERIAL PRIMARY KEY,
    user_id   BIGINT NOT NULL,

    title     VARCHAR(256),

    -- где лежит объект
    bucket_name TEXT NOT NULL,
    object_key  TEXT NOT NULL,

    -- базовая мета
    size_bytes  BIGINT NOT NULL CHECK (size_bytes >= 0),
    content_type TEXT,
    etag         TEXT,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_file_data_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE CASCADE,

    CONSTRAINT uq_file_object UNIQUE (bucket_name, object_key)
);

CREATE TABLE bank_data (
    id        BIGSERIAL PRIMARY KEY,
    user_id   BIGINT NOT NULL,

    bank_name TEXT,
    pid       BYTEA,

    CONSTRAINT fk_bank_data_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

-- оставляю название как в схеме: accout_data (если это опечатка и хочешь account_data — скажи)
CREATE TABLE account_data (
    id       BIGSERIAL PRIMARY KEY,
    user_id  BIGINT NOT NULL,

	service_name TEXT,
    username TEXT,
    password BYTEA,

    CONSTRAINT fk_accout_data_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

CREATE TABLE user_tokens (
    id                       BIGSERIAL PRIMARY KEY,
    user_id                  BIGINT NOT NULL,

    refresh_token            TEXT NOT NULL,      -- hash only
    refresh_token_expires_at timestamptz NOT NULL,    -- epoch seconds/ms (как у тебя в схеме)
    revoked_at               timestamptz,             -- epoch seconds/ms

    CONSTRAINT fk_user_tokens_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

`

// fixme: NOT GUD, setup database bu sql script in container ?

func SetupDatabase(db *sql.DB) error {
	_, err := db.Exec(createTable)
	return err
}
