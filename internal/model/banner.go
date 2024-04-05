package model

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

func (b *Banner) Validate() bool {
	if !b.BannerItem.Validate() || b.TagIds == nil || b.FeatureId <= 0 {
		return false
	}

	return true
}

func (b *BannerItem) Validate() bool {
	if b.Text == "" || b.Title == "" || b.Url == "" {
		return false
	}

	return true
}
