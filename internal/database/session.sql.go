package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func (q *Queries) GetAllSession(ctx context.Context) ([]Session, error) {
	rows, err := q.db.QueryContext(ctx, "SELECT * FROM sessions;")
	if err != nil {
		return nil, err
	}
	var s []Session
	for rows.Next() {
		var i Session
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.SessionID,
			&i.IsAuthenticated,
			&i.ExpiredAt,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		s = append(s, i)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return s, nil
}

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

const deleteSessionBySessionID = `DELETE FROM sessions WHERE session_id=?`

func (q *Queries) DeleteSessionBySessionID(ctx context.Context, sessionID string) error {
	_, err := q.db.ExecContext(ctx, deleteSessionBySessionID, sessionID)
	return err
}

type CreateSessionParams struct {
	SessionID string
	ExpiredAt int64
	UserID    sql.NullString
}

const createSession = `INSERT INTO sessions (id, created_at, updated_at, session_id, expired_at, user_id)
VALUES (?, ?, ?, ?, ?, ?);`

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) error {
	uid := uuid.New().String()
	_, err := q.db.ExecContext(ctx, createSession,
		uid,
		time.Now().Local().String(),
		time.Now().Local().String(),
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

const updateSessionAuthentication = `UPDATE sessions
SET updated_at = ?, is_authenticated = 1, user_id = ? WHERE session_id = ?;`

func (q *Queries) UpdateSessionAuthentication(ctx context.Context, arg UpdateSessionParams) error {
	_, err := q.db.ExecContext(ctx, updateSessionAuthentication,
		time.Now().Local().String(),
		arg.UserID,
		arg.SessionID,
	)
	return err
}

func (q *Queries) DeleteExpiredSession(ctx context.Context) error {
	nowInt := time.Now().Local().UnixMilli()
	_, err := q.db.ExecContext(ctx, "DELETE FROM sessions WHERE expired_at < ?;", nowInt)
	return err
}
