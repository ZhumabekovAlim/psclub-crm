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
       INSERT INTO users (name, phone, password, company_id, branch_id, role, permissions, salary_hookah, hookah_salary_type, salary_bar, salary_shift, created_at, updated_at)
       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
	res, err := r.db.ExecContext(ctx, query,
		u.Name, u.Phone, u.Password, u.CompanyID, u.BranchID, u.Role, permissionsJSON,
		u.SalaryHookah, u.HookahSalaryType, u.SalaryBar, u.SalaryShift)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

func (r *UserRepository) GetAll(ctx context.Context, companyID, branchID int) ([]models.User, error) {
	query := `SELECT id, name, phone, password, company_id, branch_id, role, permissions, salary_hookah, hookah_salary_type, salary_bar, salary_shift, created_at, updated_at FROM users WHERE role != 'director' AND company_id=? AND branch_id=? ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query, companyID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var permStr sql.NullString

		err := rows.Scan(&u.ID, &u.Name, &u.Phone, &u.Password, &u.CompanyID, &u.BranchID, &u.Role, &permStr,
			&u.SalaryHookah, &u.HookahSalaryType, &u.SalaryBar, &u.SalaryShift, &u.CreatedAt, &u.UpdatedAt)
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

func (r *UserRepository) GetByID(ctx context.Context, id, companyID, branchID int) (*models.User, error) {
	query := `SELECT id, name, phone, password, company_id, branch_id, role, permissions, salary_hookah, hookah_salary_type, salary_bar, salary_shift, created_at, updated_at FROM users WHERE id=? AND company_id=? AND branch_id=?`

	var u models.User
	var permStr sql.NullString

	err := r.db.QueryRowContext(ctx, query, id, companyID, branchID).Scan(
		&u.ID, &u.Name, &u.Phone, &u.Password, &u.CompanyID, &u.BranchID, &u.Role, &permStr,
		&u.SalaryHookah, &u.HookahSalaryType, &u.SalaryBar, &u.SalaryShift, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if permStr.Valid {
		_ = json.Unmarshal([]byte(permStr.String), &u.Permissions)
	}

	return &u, nil
}

// GetByIDNoTenant fetches user by id without tenant filtering.
func (r *UserRepository) GetByIDNoTenant(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name, phone, password, company_id, branch_id, role, permissions, salary_hookah, hookah_salary_type, salary_bar, salary_shift, created_at, updated_at FROM users WHERE id=?`

	var u models.User
	var permStr sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Name, &u.Phone, &u.Password, &u.CompanyID, &u.BranchID, &u.Role, &permStr,
		&u.SalaryHookah, &u.HookahSalaryType, &u.SalaryBar, &u.SalaryShift, &u.CreatedAt, &u.UpdatedAt)
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
                      SET name=?, phone=?, company_id=?, branch_id=?, role=?, permissions=?,
                              salary_hookah=?, hookah_salary_type=?, salary_bar=?, salary_shift=?, updated_at=NOW()
                      WHERE id=? AND company_id=? AND branch_id=?`
		args = []interface{}{
			u.Name, u.Phone, u.CompanyID, u.BranchID, u.Role, permissionsJSON,
			u.SalaryHookah, u.HookahSalaryType, u.SalaryBar, u.SalaryShift, u.ID, u.CompanyID, u.BranchID,
		}
	} else {

		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
		if err != nil {
			return err
		}
		u.Password = string(hashed)

		query = `
                      UPDATE users
                      SET name=?, phone=?, password=?, company_id=?, branch_id=?, role=?, permissions=?,
                              salary_hookah=?, hookah_salary_type=?, salary_bar=?, salary_shift=?, updated_at=NOW()
                      WHERE id=? AND company_id=? AND branch_id=?`
		args = []interface{}{
			u.Name, u.Phone, u.Password, u.CompanyID, u.BranchID, u.Role, permissionsJSON,
			u.SalaryHookah, u.HookahSalaryType, u.SalaryBar, u.SalaryShift, u.ID, u.CompanyID, u.BranchID,
		}
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id, companyID, branchID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id=? AND company_id=? AND branch_id=?`, id, companyID, branchID)
	return err
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := `SELECT id, name, phone, password, company_id, branch_id, role, permissions, salary_hookah, hookah_salary_type, salary_bar, salary_shift, created_at, updated_at FROM users WHERE phone=?`

	var u models.User
	var permStr sql.NullString

	err := r.db.QueryRowContext(ctx, query, phone).Scan(
		&u.ID, &u.Name, &u.Phone, &u.Password, &u.CompanyID, &u.BranchID, &u.Role, &permStr,
		&u.SalaryHookah, &u.HookahSalaryType, &u.SalaryBar, &u.SalaryShift, &u.CreatedAt, &u.UpdatedAt)
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
