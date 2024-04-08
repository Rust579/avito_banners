package postgres

import (
	"avito_banners/internal/model"
	"encoding/json"
	"strconv"
)

func GetAllBanners() ([]model.Banner, error) {

	var bannerData []byte
	var tagIdsData []byte

	query := "SELECT banner_id, feature_id, tag_ids, banner_data, is_active, created_at, updated_at FROM banners"

	rows, err := psgDb.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []model.Banner

	for rows.Next() {
		var banner model.Banner

		if err := rows.Scan(&banner.BannerId, &banner.FeatureId, &tagIdsData, &bannerData, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
			return nil, err
		}

		banner.CreatedAt = banner.CreatedAt.UTC()
		banner.UpdatedAt = banner.UpdatedAt.UTC()

		if err := json.Unmarshal(bannerData, &banner.BannerItem); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(tagIdsData, &banner.TagIds); err != nil {
			return nil, err
		}

		banners = append(banners, banner)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

func InsertBanner(banner model.Banner) (int, error) {

	bannerItemJson, err := json.Marshal(banner.BannerItem)
	if err != nil {
		return 0, err
	}

	tagIdsJson, err := json.Marshal(banner.TagIds)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO banners (feature_id, tag_ids, banner_data, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING banner_id`

	var bannerID int

	err = psgDb.QueryRow(query, banner.FeatureId, tagIdsJson, bannerItemJson, banner.IsActive, banner.CreatedAt, banner.UpdatedAt).Scan(&bannerID)
	if err != nil {
		return 0, err
	}

	return bannerID, nil
}

func FindBannerByFeatureAndTagId(input model.BannerGetRequest) (*model.Banner, error) {

	query := "SELECT * FROM banners " +
		"WHERE feature_id = " + strconv.Itoa(input.FeatureId) + " " +
		"AND tag_ids @> '[" + strconv.Itoa(input.TagId) + "]'::jsonb " +
		"ORDER BY updated_at DESC " +
		"LIMIT 1"

	var bannerData []byte
	var tagIdsData []byte

	row := psgDb.QueryRow(query)

	var banner model.Banner

	if err := row.Scan(&banner.BannerId, &banner.FeatureId, &tagIdsData, &bannerData, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
		return nil, err
	}

	banner.CreatedAt = banner.CreatedAt.UTC()
	banner.UpdatedAt = banner.UpdatedAt.UTC()

	if err := json.Unmarshal(bannerData, &banner.BannerItem); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(tagIdsData, &banner.TagIds); err != nil {
		return nil, err
	}

	if err := row.Err(); err != nil {
		return nil, err
	}

	return &banner, nil
}

func DeleteBannerByID(bannerID int) error {
	query := "DELETE FROM banners WHERE banner_id = $1"

	_, err := psgDb.Exec(query, bannerID)
	if err != nil {
		return err
	}

	return nil
}
