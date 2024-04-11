package pulls

import (
	"avito_banners/internal/model"
	"reflect"
	"sync"
)

// Кэш баннеров с ключами = banner_id
var bannerByIdPull = &bannerIdPullMap{pull: make(map[int]model.Banner)}

// Кэш баннеров с ключами = feature_id
var bannerByFeatureIdPull = &bannerFIdPullMap{pull: make(map[int][]model.Banner)}

type bannerIdPullMap struct {
	sync.Mutex
	pull map[int]model.Banner
}

type bannerFIdPullMap struct {
	sync.Mutex
	pull map[int][]model.Banner
}

// InitBannersPulls Сборка всего кэша
// bannerByIdPull хранит по одному баннеру на ключ
// bannerByFeatureIdPull может хранить по несколько баннеров на ключ
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
}

// AddPullBanner Добавление баннера
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
}

// UpdatePullBanner Обновления баннера
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
}

// GetBannersByFeatureId Получение баннеров по feature_id
func GetBannersByFeatureId(fId int) []model.Banner {
	bannerByFeatureIdPull.Lock()
	defer bannerByFeatureIdPull.Unlock()

	banners, ok := bannerByFeatureIdPull.pull[fId]
	if !ok {
		return nil
	}

	return banners
}

// GetBannerById Получение баннера по banner_id
func GetBannerById(id int) *model.Banner {
	bannerByIdPull.Lock()
	defer bannerByIdPull.Unlock()

	banner, ok := bannerByIdPull.pull[id]
	if !ok {
		return nil
	}

	return &banner
}

// DeleteBannerById Удаление баннера по banner_id
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
