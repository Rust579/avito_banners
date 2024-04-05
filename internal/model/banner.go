package model

import "avito_banners/internal/errs"

type Banner struct {
	TagIds     []int      `json:"tag_ids"`
	FeatureId  int        `json:"feature_id"`
	BannerItem BannerItem `json:"banner_item"`
}

type BannerItem struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

func (b *Banner) Validate() (ers []errs.Error) {
	ers = b.BannerItem.Validate()

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

func (b *BannerItem) Validate() (ers []errs.Error) {
	if b.Text == "" {
		ers = append(ers, errs.GetErr(101))
	}

	if b.Title == "" {
		ers = append(ers, errs.GetErr(102))
	}

	if b.Url == "" {
		ers = append(ers, errs.GetErr(103))
	}

	return
}
