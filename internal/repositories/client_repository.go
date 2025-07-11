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
        INSERT INTO clients (name, phone, date_of_birth, channel_id, bonus, visits, income, status, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
        res, err := r.db.ExecContext(ctx, query, c.Name, c.Phone, c.DateOfBirth, c.ChannelID, c.Bonus, c.Visits, c.Income, c.Status)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *ClientRepository) GetAll(ctx context.Context) ([]models.Client, error) {
        query := `SELECT c.id, c.name, c.phone, c.date_of_birth, c.channel_id, IFNULL(ch.name, ''), c.bonus, c.visits, c.income, c.status, c.created_at, c.updated_at
                FROM clients c
                LEFT JOIN channels ch ON c.channel_id = ch.id
                WHERE c.status <> 'deleted'
                ORDER BY c.id`
        rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.Client
	for rows.Next() {
		var c models.Client
		var dob sql.NullTime
                err := rows.Scan(&c.ID, &c.Name, &c.Phone, &dob, &c.ChannelID, &c.Channel, &c.Bonus, &c.Visits, &c.Income, &c.Status, &c.CreatedAt, &c.UpdatedAt)
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
        query := `SELECT c.id, c.name, c.phone, c.date_of_birth, c.channel_id, IFNULL(ch.name, ''), c.bonus, c.visits, c.income, c.status, c.created_at, c.updated_at
                FROM clients c
                LEFT JOIN channels ch ON c.channel_id = ch.id
                WHERE c.id=?`
	var c models.Client
	var dob sql.NullTime
        err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.Phone, &dob, &c.ChannelID, &c.Channel, &c.Bonus, &c.Visits, &c.Income, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if dob.Valid {
		c.DateOfBirth = &dob.Time
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClientRepository) Update(ctx context.Context, c *models.Client) error {
	// Получить текущие значения клиента
	existing, err := r.GetByID(ctx, c.ID)
	if err != nil {
		return err
	}

	// Прибавить к существующим значениям
	c.Bonus += existing.Bonus
	c.Visits += existing.Visits
	c.Income += existing.Income

	// Обновить поля клиента
	query := `
        UPDATE clients
        SET name = ?,
            phone = ?,
            date_of_birth = ?,
            channel_id = ?,
            bonus = ?,
            visits = ?,
            income = ?,
            updated_at = NOW()
        WHERE id = ?`

	_, err = r.db.ExecContext(ctx, query, c.Name, c.Phone, c.DateOfBirth, c.ChannelID, c.Bonus, c.Visits, c.Income, c.ID)
	return err
}

func (r *ClientRepository) Delete(ctx context.Context, id int) error {
        _, err := r.db.ExecContext(ctx, `UPDATE clients SET status='deleted' WHERE id=?`, id)
        return err
}

func (r *ClientRepository) AddBonus(ctx context.Context, clientID int, bonus int) error {
	query := `UPDATE clients SET bonus = bonus + ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, bonus, clientID)
	return err
}

func (r *ClientRepository) GetByPhone(ctx context.Context, phone string) (*models.Client, error) {
        query := `SELECT c.id, c.name, c.phone, c.date_of_birth, c.channel_id, IFNULL(ch.name, ''), c.bonus, c.visits, c.income, c.status, c.created_at, c.updated_at
                FROM clients c
                LEFT JOIN channels ch ON c.channel_id = ch.id
                WHERE c.phone = ?`
	var c models.Client
	var dob sql.NullTime
        err := r.db.QueryRowContext(ctx, query, phone).Scan(&c.ID, &c.Name, &c.Phone, &dob, &c.ChannelID, &c.Channel, &c.Bonus, &c.Visits, &c.Income, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if dob.Valid {
		c.DateOfBirth = &dob.Time
	}
	return &c, nil
}

func (r *ClientRepository) AddVisits(ctx context.Context, clientID int, visits int) error {
	query := `UPDATE clients SET visits = visits + ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, visits, clientID)
	return err
}

func (r *ClientRepository) AddIncome(ctx context.Context, clientID int, income int) error {
	query := `UPDATE clients SET income = income + ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, income, clientID)
	return err
}
