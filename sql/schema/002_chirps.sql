-- +goose Up
CREATE TABLE chirps (
    id uuid primary key not null,
    body text not null,
    user_id uuid not null,
    created_at timestamp not null,
    updated_at timestamp not null,

    CONSTRAINT chirp_user_id_fk FOREIGN KEY (user_id)
        REFERENCES users (id) on DELETE cascade
);

-- +goose Down
DROP TABLE chirps;
