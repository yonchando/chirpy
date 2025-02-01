-- +goose Up
CREATE TABLE users (
    id uuid primary key not null,
    email varchar(150) unique not null,
    hashed_password varchar(255) not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
DROP TABLE users;
