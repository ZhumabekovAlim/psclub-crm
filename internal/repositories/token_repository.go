package repositories

import (
	"context"
	"database/sql"
	"time"
)

type TokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Save(ctx context.Context, userID int, hash string, exp time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at) VALUES (?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query, userID, hash, exp)
	return err
}

func (r *TokenRepository) Delete(ctx context.Context, hash string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE token_hash=?`, hash)
	return err
}

func (r *TokenRepository) Get(ctx context.Context, hash string) (int, time.Time, error) {
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash=?`
	var userID int
	var exp time.Time
	err := r.db.QueryRowContext(ctx, query, hash).Scan(&userID, &exp)
	if err != nil {
		return 0, time.Time{}, err
	}
	return userID, exp, nil
}

func (r *TokenRepository) DeleteByUser(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id=?`, userID)
	return err
}
