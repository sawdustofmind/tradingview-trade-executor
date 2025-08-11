package repository

import (
	"context"

	"github.com/frenswifbenefits/myfren/internal/entity"
)

func (repo *Repository) GetPortfolios() ([]entity.Portfolio, error) {
	models := make([]entity.Portfolio, 0)
	err := repo.db.NewSelect().
		Model(&models).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (repo *Repository) GetPortfolioById(id int64) (entity.Portfolio, error) {
	model := entity.Portfolio{}
	err := repo.db.NewSelect().
		Model(&model).
		Where("id = ?", id).
		Scan(context.Background(), &model)
	if err != nil {
		return model, err
	}
	return model, nil
}

func (repo *Repository) FindPortfolio(name string) (entity.Portfolio, error) {
	model := entity.Portfolio{}
	err := repo.db.NewSelect().
		Model(&model).
		Where("name = ?", name).
		Scan(context.Background(), &model)
	if err != nil {
		return model, err
	}
	return model, nil
}

func (repo *Repository) InsertPortfolio(portfolio entity.Portfolio) (int64, error) {
	var id int64
	err := repo.db.NewInsert().
		Model(&portfolio).
		Returning("id").
		Scan(context.Background(), &id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) DeletePortfolio(id int64) error {
	_, err := repo.db.NewDelete().
		Model(&entity.Portfolio{Id: id}).
		WherePK().
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) UpdatePortfolio(portfolio entity.Portfolio) error {
	_, err := repo.db.NewUpdate().
		Model(&portfolio).
		WherePK().
		Exec(context.Background())

	if err != nil {
		return err
	}
	return nil
}
