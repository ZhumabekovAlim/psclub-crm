package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"psclub-crm/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) (int, error) {
	permissionsJSON, err := json.Marshal(u.Permissions)
	if err != nil {
		return 0, err
	}

	query := `
        INSERT INTO users (name, phone, password, role, permissions, salary_hookah, salary_bar, salary_shift, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
	res, err := r.db.ExecContext(ctx, query,
		u.Name, u.Phone, u.Password, u.Role, permissionsJSON,
		u.SalaryHookah, u.SalaryBar, u.SalaryShift)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, name, phone, password, role, permissions, salary_hookah, salary_bar, salary_shift, created_at, updated_at FROM users WHERE role != 'director' ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var permStr sql.NullString

		err := rows.Scan(&u.ID, &u.Name, &u.Phone, &u.Password, &u.Role, &permStr,
			&u.SalaryHookah, &u.SalaryBar, &u.SalaryShift, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if permStr.Valid {
			_ = json.Unmarshal([]byte(permStr.String), &u.Permissions)
		}

		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name, phone, password, role, permissions, salary_hookah, salary_bar, salary_shift, created_at, updated_at FROM users WHERE id=?`

	var u models.User
	var permStr sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.Phone, &u.Password, &u.Role, &permStr,
		&u.SalaryHookah, &u.SalaryBar, &u.SalaryShift, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if permStr.Valid {
		_ = json.Unmarshal([]byte(permStr.String), &u.Permissions)
	}

	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *models.User) error {
	permissionsJSON, err := json.Marshal(u.Permissions)
	if err != nil {
		return err
	}

	var query string
	var args []interface{}

	if u.Password == "" {
		query = `
			UPDATE users 
			SET name=?, phone=?, role=?, permissions=?, 
				salary_hookah=?, salary_bar=?, salary_shift=?, updated_at=NOW() 
			WHERE id=?`
		args = []interface{}{
			u.Name, u.Phone, u.Role, permissionsJSON,
			u.SalaryHookah, u.SalaryBar, u.SalaryShift, u.ID,
		}
	} else {

		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
		if err != nil {
			return err
		}
		u.Password = string(hashed)

		query = `
			UPDATE users 
			SET name=?, phone=?, password=?, role=?, permissions=?, 
				salary_hookah=?, salary_bar=?, salary_shift=?, updated_at=NOW() 
			WHERE id=?`
		args = []interface{}{
			u.Name, u.Phone, u.Password, u.Role, permissionsJSON,
			u.SalaryHookah, u.SalaryBar, u.SalaryShift, u.ID,
		}
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id=?`, id)
	return err
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := `SELECT id, name, phone, password, role, permissions, salary_hookah, salary_bar, salary_shift, created_at, updated_at FROM users WHERE phone=?`

	var u models.User
	var permStr sql.NullString

	err := r.db.QueryRowContext(ctx, query, phone).Scan(
		&u.ID, &u.Name, &u.Phone, &u.Password, &u.Role, &permStr,
		&u.SalaryHookah, &u.SalaryBar, &u.SalaryShift, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if permStr.Valid {
		_ = json.Unmarshal([]byte(permStr.String), &u.Permissions)
	}

	return &u, nil
}
