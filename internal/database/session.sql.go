package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const getSessionByID = `SELECT * FROM sessions WHERE session_id=?`

func (q *Queries) GetSessionBySessionID(ctx context.Context, sessionID string) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSessionByID, sessionID)
	var i Session
	if err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SessionID,
		&i.IsAuthenticated,
		&i.ExpiredAt,
		&i.UserID,
	); err != nil {
		return i, err
	}
	return i, nil
}

const deleteSessionBySessionID = `DELETE sessions WHERE session_id=?`

func (q *Queries) DeleteSessionBySessionID(ctx context.Context, sessionID string) error {
	_, err := q.db.ExecContext(ctx, deleteSessionBySessionID, sessionID)
	return err
}

type CreateSessionParams struct {
	SessionID string
	ExpiredAt int64
	UserID    sql.NullString
}

const createSession = `INSERT INTO sessions (id, session_id, expired_at, user_id)
VALUES (?, ?, ?, ?);`

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) error {
	uid := uuid.New().String()
	_, err := q.db.ExecContext(ctx, createSession,
		uid,
		arg.SessionID,
		arg.ExpiredAt,
		arg.UserID,
	)
	return err
}

type UpdateSessionParams struct {
	SessionID string
	UserID    sql.NullString
}

const updateSessionAuthentication = `UPDATE TABLE sessions
SET is_authenticated = 1, user_id = ? WHERE session_id = ?`

func (q *Queries) UpdateSessionAuthentication(ctx context.Context, arg UpdateSessionParams) error {
	_, err := q.db.ExecContext(ctx, updateSessionAuthentication, arg.UserID, arg.SessionID)
	return err
}
