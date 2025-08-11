package repository

import (
	"context"

	"github.com/frenswifbenefits/myfren/internal/entity"
)

func (repo *Repository) GetCustomerByUsername(username string) (entity.Customer, error) {
	model := entity.Customer{}
	err := repo.db.NewSelect().
		Model(&model).
		Where("LOWER(username) = LOWER(?)", username).
		Scan(context.Background())
	if err != nil {
		return model, err
	}
	return model, nil
}

func (repo *Repository) GetCustomers() ([]entity.Customer, error) {
	models := make([]entity.Customer, 0)
	err := repo.db.NewSelect().
		Model(&models).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (repo *Repository) InsertCustomer(customer entity.Customer) (int64, error) {
	var id int64
	err := repo.db.NewInsert().
		Model(&customer).
		Returning("id").
		Scan(context.Background(), &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *Repository) UpdateCustomer(customer entity.Customer) error {
	_, err := repo.db.NewUpdate().
		Model(&customer).
		Set("bybit_api_key = ?", customer.BybitApiKey).
		Set("bybit_test_api_key = ?", customer.BybitTestApiKey).
		Set("bybit_api_secret = ?", customer.BybitApiSecret).
		Set("bybit_test_api_secret = ?", customer.BybitTestApiSecret).
		Set("legal_name = ?", customer.LegalName).
		Set("gender = ?", customer.Gender).
		Set("country = ?", customer.Country).
		Set("phone_number = ?", customer.PhoneNumber).
		WherePK().
		Exec(context.Background())

	if err != nil {
		return err
	}
	return nil
}
