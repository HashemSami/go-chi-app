```go
// 2.1 Query for user's session
	// 2.2 If found, update the user's session
	// 2.3 if not found, create a new session for the user
	row := ss.DB.QueryRow(`
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
		RETURNING id;
		`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)

	if err == sql.ErrNoRows {
		row = ss.DB.QueryRow(`
		INSERT INTO sessions (user_id, token_hash)
		VALUES ($1,$2)
		RETURNING id;
		`, session.UserID, session.TokenHash)
		err = row.Scan(&session.ID)
	}
```
