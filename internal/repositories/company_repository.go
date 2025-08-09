package repositories

import (
	"context"
	"database/sql"
)

type CompanyRepository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) CreateCompany(ctx context.Context, name string) (int, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO companies (name) VALUES (?)`, name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *CompanyRepository) CreateBranch(ctx context.Context, companyID int, name string) (int, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO branches (company_id, name) VALUES (?, ?)`, companyID, name)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}
