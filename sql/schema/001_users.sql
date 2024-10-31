
-- +goose Up
CREATE TABLE Users (
id uuid PRIMARY KEY,
created_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP,
email Text unique not null
);

-- +goose Down
DROP TABLE users;
