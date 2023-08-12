-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions(
  id SERIAL PRIMARY KEY,
  -- ON DELETE CASCADE deletes the row if the related
  -- table row got deleted
  user_id INT UNIQUE REFERENCES users (id) ON DELETE CASCADE,
  token_hash TEXT UNIQUE NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
