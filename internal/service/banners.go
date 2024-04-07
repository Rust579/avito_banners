package service

import (
	"avito_banners/internal/errs"
	"avito_banners/internal/model"
	"avito_banners/internal/repo/postgres"
	"avito_banners/internal/response"
	"avito_banners/internal/service/pulls"
	"reflect"
	"time"
)

func AddBanner(bannerData model.BannerAddRequest, resp *response.Response) int {

	banner := model.Banner{
		FeatureId:  bannerData.FeatureId,
		TagIds:     bannerData.TagIds,
		BannerItem: bannerData.BannerItem,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	bId, ok := checkExistsBanner(banner)
	if ok {
		resp.SetError(errs.GetErr(111))
		resp.SetValue(bId)
		return 0
	}

	id, err := postgres.InsertBanner(banner)
	if err != nil {
		resp.SetError(errs.GetErr(99, err.Error()))
		return 0
	}

	banner.BannerId = id
	pulls.AddPullBanner(banner)

	return id
}

func checkExistsBanner(data model.Banner) (int, bool) {

	banners := pulls.GetBannerBuFeatureId(data.FeatureId)
	if len(banners) == 0 {
		return 0, false
	}

	for _, b := range banners {
		if reflect.DeepEqual(data.TagIds, b.TagIds) {
			return b.BannerId, true
		}
	}

	return 0, false
}
