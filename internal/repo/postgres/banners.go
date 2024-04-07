package postgres

import (
	"avito_banners/internal/model"
	"encoding/json"
)

func GetAllBanners() ([]model.Banner, error) {

	var bannerData []byte
	var tagIdsData []byte

	query := "SELECT banner_id, feature_id, tag_ids, banner_data, created_at, updated_at FROM banners"

	rows, err := psgDb.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []model.Banner

	for rows.Next() {
		var banner model.Banner

		if err := rows.Scan(&banner.BannerId, &banner.FeatureId, &tagIdsData, &bannerData, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
			return nil, err
		}

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

	query := `INSERT INTO banners (feature_id, tag_ids, banner_data, created_at, updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING banner_id`

	var bannerID int

	err = psgDb.QueryRow(query, banner.FeatureId, tagIdsJson, bannerItemJson, banner.CreatedAt, banner.UpdatedAt).Scan(&bannerID)
	if err != nil {
		return 0, err
	}

	return bannerID, nil
}
