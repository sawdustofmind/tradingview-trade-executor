package repository

import (
	"context"

	"github.com/frenswifbenefits/myfren/internal/entity"
)

func (repo *Repository) GetActionTrades(customerId int64) ([]entity.ActionTrade, error) {
	var models []entity.ActionTrade
	err := repo.db.NewSelect().
		Model(&models).
		Where("customer_id = ?", customerId).
		Scan(context.Background(), &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (repo *Repository) InsertActionTrade(at entity.ActionTrade) (int64, error) {
	var id int64
	err := repo.db.NewInsert().
		Model(&at).
		Returning("id").
		Scan(context.Background(), &id)
	if err != nil {
		return 0, err
	}
	return id, nil

}
