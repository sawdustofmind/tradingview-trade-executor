package repository

import (
	"context"

	"github.com/frenswifbenefits/myfren/internal/entity"
)

func (repo *Repository) GetPortfolioSubscriptionById(id int64) (entity.PortfolioSubscription, error) {
	model := entity.PortfolioSubscription{}
	err := repo.db.NewSelect().
		Model(&model).
		Where("id = ?", id).
		Scan(context.Background(), &model)
	if err != nil {
		return model, err
	}
	return model, nil
}

func (repo *Repository) GetPortfolioSubscription(customerId int64) ([]entity.PortfolioSubscription, error) {
	models := make([]entity.PortfolioSubscription, 0)
	err := repo.db.NewSelect().
		Model(&models).
		Where("customer_id = ?", customerId).
		Scan(context.Background(), &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (repo *Repository) GetAllActiveSubscription() ([]entity.PortfolioSubscription, error) {
	models := make([]entity.PortfolioSubscription, 0)
	err := repo.db.NewSelect().
		Model(&models).
		Where("status = 'active'").
		Scan(context.Background(), &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (repo *Repository) GetActivePortfolioSubscriptions(portfolioId int64, exchange string) ([]entity.PortfolioSubscription, error) {
	models := make([]entity.PortfolioSubscription, 0)
	err := repo.db.NewSelect().
		Model(&models).
		Where("status = 'active'").
		Where("portfolio_id = ?", portfolioId).
		Where("exchange = ?", exchange).
		Scan(context.Background(), &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (repo *Repository) InsertPortfolioSubscription(subscription entity.PortfolioSubscription) (int64, error) {
	var id int64
	err := repo.db.NewInsert().
		Model(&subscription).
		Returning("id").
		Scan(context.Background(), &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *Repository) PortfolioUnsubscribe(customerId int64, subscriptionId int64) error {
	_, err := repo.db.NewUpdate().
		Model(&entity.PortfolioSubscription{}).
		Set("status = ?", "terminated").
		Where("customer_id = ?", customerId).
		Where("id = ?", subscriptionId).
		Exec(context.Background())

	if err != nil {
		return err
	}
	return nil
}
