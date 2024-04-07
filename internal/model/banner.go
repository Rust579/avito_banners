package model

import (
	"avito_banners/internal/errs"
	"time"
)

type Banner struct {
	BannerId   int                    `json:"banner_id"`
	TagIds     []int                  `json:"tag_ids"`
	FeatureId  int                    `json:"feature_id"`
	BannerItem map[string]interface{} `json:"banner_item"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type BannerAddRequest struct {
	TagIds     []int                  `json:"tag_ids"`
	FeatureId  int                    `json:"feature_id"`
	BannerItem map[string]interface{} `json:"banner_item"`
}

func (b *BannerAddRequest) Validate() (ers []errs.Error) {

	if len(b.TagIds) == 0 {
		ers = append(ers, errs.GetErr(105))
	}

	if b.FeatureId <= 0 {
		ers = append(ers, errs.GetErr(106))
	}

	for _, t := range b.TagIds {
		if t <= 0 {
			ers = append(ers, errs.GetErr(107))
		}
	}

	return
}
