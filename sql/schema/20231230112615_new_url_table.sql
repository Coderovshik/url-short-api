-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS url(
    id SERIAL PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE url;
-- +goose StatementEnd
