-- +goose Up
CREATE TABLE refresh_tokens
(
    token      text primary key not null,
    user_id    uuid             not null,
    expires_at timestamp        not null,
    revoked_at timestamp        default null,
    created_at timestamp        not null,
    updated_at timestamp        not null,

    CONSTRAINT refresh_tokens_user_id_foreign FOREIGN KEY (user_id) REFERENCES users (id) on delete cascade
);

-- +goose Down
DROP TABLE refresh_tokens;
