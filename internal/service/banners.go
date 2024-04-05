package service

import (
	"avito_banners/internal/model"
	"avito_banners/internal/repo/postgres"
)

func AddBanner(banner model.Banner) (int, error) {

	id, err := postgres.InsertBanner(banner)
	if err != nil {
		return 0, err
	}

	return id, nil
}
