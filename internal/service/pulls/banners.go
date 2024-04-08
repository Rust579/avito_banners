package pulls

import (
	"avito_banners/internal/model"
	"sync"
)

var bannerByIdPull = &bannerIdPullMap{pull: make(map[int]model.Banner)}
var bannerByFeatureIdPull = &bannerFIdPullMap{pull: make(map[int][]model.Banner)}

type bannerIdPullMap struct {
	sync.Mutex
	pull map[int]model.Banner
}

type bannerFIdPullMap struct {
	sync.Mutex
	pull map[int][]model.Banner
}

func InitBannersPulls(banners []model.Banner) {
	for _, b := range banners {
		bannerByIdPull.Lock()
		bannerByIdPull.pull[b.BannerId] = b
		bannerByIdPull.Unlock()

		bannerByFeatureIdPull.Lock()
		bs := bannerByFeatureIdPull.pull[b.FeatureId]
		bs = append(bs, b)
		bannerByFeatureIdPull.pull[b.FeatureId] = bs
		bannerByFeatureIdPull.Unlock()
	}

	return
}

func AddPullBanner(banner model.Banner) {
	bannerByIdPull.Lock()
	_, ok := bannerByIdPull.pull[banner.BannerId]
	if !ok {
		bannerByIdPull.pull[banner.BannerId] = banner
	}
	bannerByIdPull.Unlock()

	bannerByFeatureIdPull.Lock()
	b, okk := bannerByFeatureIdPull.pull[banner.FeatureId]
	if !okk {
		bannerByFeatureIdPull.pull[banner.FeatureId] = []model.Banner{banner}
	}

	b = append(b, banner)
	bannerByFeatureIdPull.pull[banner.FeatureId] = b

	bannerByFeatureIdPull.Unlock()

	return
}

func GetBannersByFeatureId(fId int) []model.Banner {
	bannerByFeatureIdPull.Lock()
	defer bannerByFeatureIdPull.Unlock()

	banners, ok := bannerByFeatureIdPull.pull[fId]
	if !ok {
		return nil
	}

	return banners
}

func GetBannerById(id int) *model.Banner {
	bannerByIdPull.Lock()
	defer bannerByIdPull.Unlock()

	banner, ok := bannerByIdPull.pull[id]
	if !ok {
		return nil
	}

	return &banner
}
