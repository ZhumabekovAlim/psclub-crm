package repositories

import (
	"context"
	"database/sql"
	"psclub-crm/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) (int, error) {
	query := `
        INSERT INTO users (name, phone, password, role, salary_hookah, salary_bar, salary_shift, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
	res, err := r.db.ExecContext(ctx, query, u.Name, u.Phone, u.Password, u.Role, u.SalaryHookah, u.SalaryBar, u.SalaryShift)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, name, phone, password, role, salary_hookah, salary_bar, salary_shift, created_at, updated_at FROM users ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Name, &u.Phone, &u.Password, &u.Role, &u.SalaryHookah, &u.SalaryBar, &u.SalaryShift, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name, phone, password, role, salary_hookah, salary_bar, salary_shift, created_at, updated_at FROM users WHERE id=?`
	var u models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Phone, &u.Password, &u.Role, &u.SalaryHookah, &u.SalaryBar, &u.SalaryShift, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *models.User) error {
	query := `UPDATE users SET name=?, phone=?, password=?, role=?, salary_hookah=?, salary_bar=?, salary_shift=?, updated_at=NOW() WHERE id=?`
	_, err := r.db.ExecContext(ctx, query, u.Name, u.Phone, u.Password, u.Role, u.SalaryHookah, u.SalaryBar, u.SalaryShift, u.ID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id=?`, id)
	return err
}
