CREATE TABLE sessions(
  id SERIAL PRIMARY KEY,
  -- ON DELETE CASCADE deletes the row if the related
  -- table row got deleted
  user_id INT UNIQUE REFERENCES users (id) ON DELETE CASCADE,
  token_hash TEXT UNIQUE NOT NULL
);

INSERT INTO sessions (user_id, token_hash)
VALUES ($1,$2)
RETURNING id;


DELETE FROM sessions
WHERE token_hash = $1;

-- adding foreign keys when the table already exists
ALTER TABLE sessions
  ADD CONSTRAINT sessions_user_id_fkey
  FOREIGN KEY (user_id) REFERENCES users (id);


-- creating the index manually
CREATE INDEX sessions_token_hash_idx ON sessions (token_hash,user_id);

-- insert or update if exists (on postgres)
INSERT INTO sessions(user_id, token_hash)
VALUES (1, "some value") ON CONFLICT (user_id) DO
UPDATE
SET token_hash = "some value";