package service

import (
	"avito_banners/internal/model"
	"avito_banners/internal/repo/postgres"
)

func AddBanner(banner model.Banner) error {

	if err := postgres.InsertBanner(banner); err != nil {
		return err
	}
	return nil
}
