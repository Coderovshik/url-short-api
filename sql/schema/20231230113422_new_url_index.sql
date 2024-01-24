-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_alias;
-- +goose StatementEnd
