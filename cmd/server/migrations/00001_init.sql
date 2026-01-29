-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users (
                                     id            BIGSERIAL PRIMARY KEY,
                                     username      VARCHAR(255) UNIQUE NOT NULL,
                                     password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS text_data (
                                         id      BIGSERIAL PRIMARY KEY,
                                         user_id BIGINT NOT NULL,

                                         title   VARCHAR(256),
                                         text    TEXT,

                                         CONSTRAINT fk_text_data_user
                                             FOREIGN KEY (user_id)
                                                 REFERENCES users(id)
                                                 ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS file_data (
                                         id        BIGSERIAL PRIMARY KEY,
                                         user_id   BIGINT NOT NULL,

                                         title     VARCHAR(256),

                                         bucket_name TEXT NOT NULL,
                                         object_key  TEXT NOT NULL,

                                         size_bytes   BIGINT NOT NULL CHECK (size_bytes >= 0),
                                         content_type TEXT,
                                         etag         TEXT,

                                         created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

                                         CONSTRAINT fk_file_data_user
                                             FOREIGN KEY (user_id)
                                                 REFERENCES users(id)
                                                 ON DELETE CASCADE,

                                         CONSTRAINT uq_file_object UNIQUE (bucket_name, object_key)
);

CREATE TABLE IF NOT EXISTS bank_data (
                                         id        BIGSERIAL PRIMARY KEY,
                                         user_id   BIGINT NOT NULL,

                                         bank_name TEXT,
                                         pid       BYTEA,

                                         CONSTRAINT fk_bank_data_user
                                             FOREIGN KEY (user_id)
                                                 REFERENCES users(id)
                                                 ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS account_data (
                                            id       BIGSERIAL PRIMARY KEY,
                                            user_id  BIGINT NOT NULL,

                                            service_name TEXT,
                                            username     TEXT,
                                            password     BYTEA,

                                            CONSTRAINT fk_account_data_user
                                                FOREIGN KEY (user_id)
                                                    REFERENCES users(id)
                                                    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_tokens (
                                           id                       BIGSERIAL PRIMARY KEY,
                                           user_id                  BIGINT NOT NULL,

                                           refresh_token            TEXT NOT NULL,
                                           refresh_token_expires_at TIMESTAMPTZ NOT NULL,
                                           revoked_at               TIMESTAMPTZ,

                                           CONSTRAINT fk_user_tokens_user
                                               FOREIGN KEY (user_id)
                                                   REFERENCES users(id)
                                                   ON DELETE CASCADE
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS user_tokens CASCADE;
DROP TABLE IF EXISTS account_data CASCADE;
DROP TABLE IF EXISTS bank_data CASCADE;
DROP TABLE IF EXISTS file_data CASCADE;
DROP TABLE IF EXISTS text_data CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- +goose StatementEnd