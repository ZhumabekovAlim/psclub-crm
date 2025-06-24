package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type ChannelRepository struct {
	db *sql.DB
}

func NewChannelRepository(db *sql.DB) *ChannelRepository {
	return &ChannelRepository{db: db}
}

func (r *ChannelRepository) Create(ctx context.Context, ch *models.Channel) (int, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO channels (name) VALUES (?)`, ch.Name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *ChannelRepository) GetAll(ctx context.Context) ([]models.Channel, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM channels ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Channel
	for rows.Next() {
		var ch models.Channel
		if err := rows.Scan(&ch.ID, &ch.Name); err != nil {
			return nil, err
		}
		result = append(result, ch)
	}
	return result, nil
}

func (r *ChannelRepository) Update(ctx context.Context, ch *models.Channel) error {
	_, err := r.db.ExecContext(ctx, `UPDATE channels SET name=? WHERE id=?`, ch.Name, ch.ID)
	return err
}

func (r *ChannelRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM channels WHERE id=?`, id)
	return err
}

func (r *ChannelRepository) GetByID(ctx context.Context, id int) (*models.Channel, error) {
	var ch models.Channel
	err := r.db.QueryRowContext(ctx, `SELECT id, name FROM channels WHERE id=?`, id).Scan(&ch.ID, &ch.Name)
	if err != nil {
		return nil, err
	}
	return &ch, nil
}
