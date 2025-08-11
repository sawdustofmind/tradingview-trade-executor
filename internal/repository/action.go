package repository

import (
	"context"

	"github.com/frenswifbenefits/myfren/internal/entity"
	"github.com/google/uuid"
)

func (repo *Repository) GetActions(customerId int64) ([]entity.Action, error) {
	var models []entity.Action
	err := repo.db.NewSelect().
		Model(&models).
		Where("customer_id = ?", customerId).
		Scan(context.Background(), &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (repo *Repository) GetUnprocessedActions() ([]entity.Action, error) {
	var models []entity.Action
	err := repo.db.NewSelect().
		Model(&models).
		Where("need_to_fetch_trades = true").
		Scan(context.Background(), &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (repo *Repository) FindAction(id uuid.UUID) (entity.Action, error) {
	var model entity.Action
	err := repo.db.NewSelect().
		Model(&model).
		Where("corr_id = ?", id.String()).
		Scan(context.Background(), &model)
	if err != nil {
		return model, err
	}
	return model, nil
}

func (repo *Repository) InsertAction(action entity.Action) (int64, error) {
	var id int64
	err := repo.db.NewInsert().
		Model(&action).
		Returning("id").
		Scan(context.Background(), &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *Repository) MarkActionAsProcessed(id int64) error {
	_, err := repo.db.NewUpdate().
		Model(&entity.Action{}).
		Set("need_to_fetch_trades = false").
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}
