package repository

import (
	"context"

	"github.com/frenswifbenefits/myfren/internal/entity"
)

func (repo *Repository) GetFrens() ([]entity.Fren, error) {
	models := make([]entity.Fren, 0)
	err := repo.db.NewSelect().
		Model(&models).
		Relation("Portfolios").
		Scan(context.Background())
	if err != nil {
		return nil, err
	}
	return models, nil

}

func (repo *Repository) GetFrenById(id int64) (entity.Fren, error) {
	model := entity.Fren{}
	err := repo.db.NewSelect().
		Model(&model).
		Relation("Portfolios").
		Where("id = ?", id).
		Scan(context.Background())
	if err != nil {
		return model, err
	}
	return model, nil
}

func (repo *Repository) InsertFren(portfolio entity.Fren) (int64, error) {
	var id int64
	err := repo.db.NewInsert().
		Model(&portfolio).
		Returning("id").
		Scan(context.Background(), &id)
	if err != nil {
		return 0, err
	}

	if len(portfolio.Portfolios) == 0 {
		return id, nil
	}

	fks := make([]entity.FrenPortfolio, 0, len(portfolio.Portfolios))
	for _, strategy := range portfolio.Portfolios {
		fks = append(fks, entity.FrenPortfolio{
			FrenID:      id,
			PortfolioId: strategy.Id,
		})
	}

	_, err = repo.db.NewInsert().
		Model(&fks).
		Exec(context.Background())
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) DeleteFren(id int64) error {
	_, err := repo.db.NewDelete().
		Model(&entity.Fren{Id: id}).
		WherePK().
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) UpdateFren(fren entity.Fren) error {
	_, err := repo.db.NewUpdate().
		Model(&fren).
		WherePK().
		Exec(context.Background())

	if err != nil {
		return err
	}
	return nil
}
