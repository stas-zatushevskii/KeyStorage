package db

import "database/sql"

var createTable = `
DROP TABLE IF EXISTS user_tokens;
DROP TABLE IF EXISTS data;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
                       id            BIGSERIAL PRIMARY KEY,
                       username      VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE data (
                      id        BIGSERIAL PRIMARY KEY,
                      user_id   BIGINT NOT NULL,

                      type      VARCHAR(255),
                      title     VARCHAR(256),
                      json_data JSONB,

                      CONSTRAINT fk_data_user
                          FOREIGN KEY (user_id)
                              REFERENCES users(id)
                              ON DELETE CASCADE
);

CREATE TABLE user_tokens (
                             id                       BIGSERIAL PRIMARY KEY,
                             user_id                  BIGINT NOT NULL,

                             refresh_token_hash       TEXT NOT NULL,
                             refresh_token_expires_at TIMESTAMPTZ NOT NULL,

                             revoked_at               TIMESTAMPTZ,

                             CONSTRAINT fk_user_tokens_user
                                 FOREIGN KEY (user_id)
                                     REFERENCES users(id)
                                     ON DELETE CASCADE
);

CREATE INDEX idx_data_user_id
    ON data(user_id);

CREATE INDEX idx_user_tokens_user_id
    ON user_tokens(user_id);

CREATE INDEX idx_user_tokens_valid
    ON user_tokens(user_id, refresh_token_expires_at)
    WHERE revoked_at IS NULL;
`

func SetupDatabase(db *sql.DB) error {
	_, err := db.Exec(createTable)
	return err
}
