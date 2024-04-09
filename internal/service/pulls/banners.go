package pulls

import (
	"avito_banners/internal/model"
	"reflect"
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

func UpdatePullBanner(updBanner model.Banner) {
	bannerByIdPull.Lock()
	_, ok := bannerByIdPull.pull[updBanner.BannerId]
	if ok {
		bannerByIdPull.pull[updBanner.BannerId] = updBanner
	}
	bannerByIdPull.Unlock()

	bannerByFeatureIdPull.Lock()
	banners, okk := bannerByFeatureIdPull.pull[updBanner.FeatureId]
	if okk {
		for i, b := range banners {
			if b.BannerId == updBanner.BannerId {
				banners[i] = updBanner
			}
		}
		bannerByFeatureIdPull.pull[updBanner.FeatureId] = banners
	}

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

func DeleteBannerById(bannerData model.Banner) {
	bannerByIdPull.Lock()
	_, ok := bannerByIdPull.pull[bannerData.BannerId]
	if ok {
		delete(bannerByIdPull.pull, bannerData.BannerId)
	}
	bannerByIdPull.Unlock()

	bannerByFeatureIdPull.Lock()
	banners, okk := bannerByFeatureIdPull.pull[bannerData.FeatureId]
	if okk {

		var bannersByVer []model.Banner
		var bannerForDelete model.Banner

		for _, b := range banners {
			if reflect.DeepEqual(bannerData.TagIds, b.TagIds) {
				bannersByVer = append(bannersByVer, b)
			}
		}

		for _, b := range bannersByVer {
			if b.UpdatedAt == bannerData.UpdatedAt {
				bannerForDelete = b
			}
		}

		for i, b := range banners {
			if b.BannerId == bannerForDelete.BannerId {
				banners = append(banners[:i], banners[i+1:]...)
			}
		}
	}

	bannerByFeatureIdPull.pull[bannerData.FeatureId] = banners
	bannerByFeatureIdPull.Unlock()
}
