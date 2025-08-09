package services

import (
	"context"
	"fmt"

	"psclub-crm/internal/common"
	"psclub-crm/internal/models"
	"psclub-crm/internal/repositories"
)

type CompanyService struct {
	repo            *repositories.CompanyRepository
	channelRepo     *repositories.ChannelRepository
	tableCatRepo    *repositories.TableCategoryRepository
	tableRepo       *repositories.TableRepository
	categoryRepo    *repositories.CategoryRepository
	paymentTypeRepo *repositories.PaymentTypeRepository
	settingsRepo    *repositories.SettingsRepository
	expenseCatRepo  *repositories.ExpenseCategoryRepository
}

func NewCompanyService(
	repo *repositories.CompanyRepository,
	channelRepo *repositories.ChannelRepository,
	tableCatRepo *repositories.TableCategoryRepository,
	tableRepo *repositories.TableRepository,
	categoryRepo *repositories.CategoryRepository,
	paymentTypeRepo *repositories.PaymentTypeRepository,
	settingsRepo *repositories.SettingsRepository,
	expenseCatRepo *repositories.ExpenseCategoryRepository,
) *CompanyService {
	return &CompanyService{
		repo:            repo,
		channelRepo:     channelRepo,
		tableCatRepo:    tableCatRepo,
		tableRepo:       tableRepo,
		categoryRepo:    categoryRepo,
		paymentTypeRepo: paymentTypeRepo,
		settingsRepo:    settingsRepo,
		expenseCatRepo:  expenseCatRepo,
	}
}

type CreateCompanyInput struct {
	Name       string
	BranchName string
}

func (s *CompanyService) CreateCompany(ctx context.Context, name, branchName string) (int, int, error) {
	companyID, err := s.repo.CreateCompany(ctx, name)
	if err != nil {
		return 0, 0, err
	}
	branchID, err := s.repo.CreateBranch(ctx, companyID, branchName)
	if err != nil {
		return 0, 0, err
	}

	tenantCtx := context.WithValue(ctx, common.CtxCompanyID, companyID)
	tenantCtx = context.WithValue(tenantCtx, common.CtxBranchID, branchID)

	if _, err := s.channelRepo.Create(tenantCtx, &models.Channel{Name: "- не указано"}); err != nil {
		return 0, 0, err
	}

	tableCatID, err := s.tableCatRepo.Create(tenantCtx, &models.TableCategory{Name: "Основной зал", CompanyID: companyID, BranchID: branchID})
	if err != nil {
		return 0, 0, err
	}
	for i := 1; i <= 6; i++ {
		if _, err := s.tableRepo.Create(tenantCtx, &models.Table{CategoryID: tableCatID, Name: fmt.Sprintf("Стол %d", i), CompanyID: companyID, BranchID: branchID}); err != nil {
			return 0, 0, err
		}
	}

	for _, c := range []string{"Бар", "Кальян", "Сеты", "Часы"} {
		if _, err := s.categoryRepo.Create(tenantCtx, &models.Category{Name: c, CompanyID: companyID, BranchID: branchID}); err != nil {
			return 0, 0, err
		}
	}

	paymentNames := []string{"Наличными", "Картой", "Каспи QR"}
	firstPTID := 0
	for i, p := range paymentNames {
		id, err := s.paymentTypeRepo.Create(tenantCtx, &models.PaymentType{Name: p, HoldPercent: 0})
		if err != nil {
			return 0, 0, err
		}
		if i == 0 {
			firstPTID = id
		}
	}

	for _, e := range []string{"Бар", "Кальян", "Ремонт", "Инвентаризация", "Зарплата", "Касса"} {
		if _, err := s.expenseCatRepo.Create(tenantCtx, &models.ExpenseCategory{Name: e}); err != nil {
			return 0, 0, err
		}
	}

	if _, err := s.settingsRepo.Create(tenantCtx, &models.Settings{
		PaymentType:      firstPTID,
		BlockTime:        120,
		BonusPercent:     0,
		WorkTimeFrom:     "10:00",
		WorkTimeTo:       "23:00",
		TablesCount:      6,
		NotificationTime: 5,
		CompanyID:        companyID,
		BranchID:         branchID,
	}); err != nil {
		return 0, 0, err
	}

	return companyID, branchID, nil
}
