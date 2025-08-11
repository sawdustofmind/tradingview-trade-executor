package repository

import (
	"context"

	"github.com/frenswifbenefits/myfren/internal/entity"
)

func (repo *Repository) GetInviteTokens() ([]string, error) {
	models := make([]entity.InviteToken, 0)
	err := repo.db.NewSelect().
		Model(&models).
		Scan(context.Background(), &models)
	if err != nil {
		return nil, err
	}
	tokens := make([]string, len(models))
	for i, m := range models {
		tokens[i] = m.Token
	}
	return tokens, nil
}

func (repo *Repository) InsertInviteToken(token string) error {
	mdl := entity.InviteToken{Token: token}
	_, err := repo.db.NewInsert().
		Model(&mdl).
		Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) DeleteInviteToken(token string) error {
	_, err := repo.db.NewDelete().
		Model(&entity.InviteToken{}).
		Where("token = ?", token).
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
