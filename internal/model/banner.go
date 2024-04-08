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
	IsActive   bool                   `json:"is_active"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type BannerAddRequest struct {
	TagIds     []int                  `json:"tag_ids"`
	FeatureId  int                    `json:"feature_id"`
	IsActive   bool                   `json:"is_active"`
	BannerItem map[string]interface{} `json:"banner_item"`
}

type BannerGetRequest struct {
	TagId           int  `json:"tag_id"`
	FeatureId       int  `json:"feature_id"`
	UseLastRevision bool `json:"use_last_revision"`
}

type BannerUpdateRequest struct {
	BannerId   int                    `json:"banner_id"`
	TagIds     []int                  `json:"tag_ids"`
	FeatureId  int                    `json:"feature_id"`
	IsActive   bool                   `json:"is_active"`
	BannerItem map[string]interface{} `json:"banner_item"`
}

type BannerVersionsRequest struct {
	TagIds    []int `json:"tag_ids"`
	FeatureId int   `json:"feature_id"`
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

func (b *BannerGetRequest) Validate() (ers []errs.Error) {

	if b.TagId <= 0 {
		ers = append(ers, errs.GetErr(106))
	}

	if b.FeatureId <= 0 {
		ers = append(ers, errs.GetErr(106))
	}

	return
}

func (b *BannerUpdateRequest) Validate() (ers []errs.Error) {

	if b.BannerId <= 0 {
		ers = append(ers, errs.GetErr(113))
	}

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
