package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

func (r *ClientRepository) Create(ctx context.Context, c *models.Client) (int, error) {
	query := `
        INSERT INTO clients (name, phone,date_of_birth, channel, bonus, visits, income, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
	res, err := r.db.ExecContext(ctx, query, c.Name, c.Phone, c.DateOfBirth, c.Channel, c.Bonus, c.Visits, c.Income)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *ClientRepository) GetAll(ctx context.Context) ([]models.Client, error) {
	query := `SELECT id, name, phone, date_of_birth, channel, bonus, visits, income, created_at, updated_at FROM clients ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Client
	for rows.Next() {
		var c models.Client
		var dob sql.NullTime
		err := rows.Scan(&c.ID, &c.Name, &c.Phone, &dob, &c.Channel, &c.Bonus, &c.Visits, &c.Income, &c.CreatedAt, &c.UpdatedAt)
		if dob.Valid {
			c.DateOfBirth = &dob.Time
		}
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *ClientRepository) GetByID(ctx context.Context, id int) (*models.Client, error) {
	query := `SELECT id, name, phone,date_of_birth, channel, bonus, visits, income, created_at, updated_at FROM clients WHERE id=?`
	var c models.Client
	var dob sql.NullTime
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.Phone, &dob, &c.Channel, &c.Bonus, &c.Visits, &c.Income, &c.CreatedAt, &c.UpdatedAt)
	if dob.Valid {
		c.DateOfBirth = &dob.Time
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClientRepository) Update(ctx context.Context, c *models.Client) error {
	query := `UPDATE clients SET name=?, phone=?, date_of_birth=?, channel=?, bonus=?, visits=?, income=?, updated_at=NOW() WHERE id=?`
	_, err := r.db.ExecContext(ctx, query, c.Name, c.Phone, c.DateOfBirth, c.Channel, c.Bonus, c.Visits, c.Income, c.ID)
	return err
}

func (r *ClientRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM clients WHERE id=?`, id)
	return err
}

func (r *ClientRepository) AddBonus(ctx context.Context, clientID int, bonus int) error {
	query := `UPDATE clients SET bonus = bonus + ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, bonus, clientID)
	return err
}
