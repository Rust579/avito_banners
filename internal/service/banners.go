package service

import (
	"avito_banners/internal/errs"
	"avito_banners/internal/model"
	"avito_banners/internal/repo/postgres"
	"avito_banners/internal/response"
	"avito_banners/internal/service/pulls"
	"reflect"
	"slices"
	"time"
)

func CreateBanner(bannerData model.BannerAddRequest, resp *response.Response) int {

	banner := model.Banner{
		FeatureId:  bannerData.FeatureId,
		TagIds:     bannerData.TagIds,
		BannerItem: bannerData.BannerItem,
		IsActive:   bannerData.IsActive,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
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

func GetBanner(bannerData model.BannerGetRequest) (*model.Banner, error) {

	var banners []model.Banner
	var bannersByVer []model.Banner

	if bannerData.UseLastRevision {

		banner, err := postgres.FindBannerByFeatureAndTagId(bannerData)
		if err != nil {
			return nil, err
		}

		return banner, nil

	} else {

		banners = pulls.GetBannersByFeatureId(bannerData.FeatureId)

		for _, b := range banners {

			slices.Sort(b.TagIds)

			_, found := slices.BinarySearch(b.TagIds, bannerData.TagId)
			if found {
				bannersByVer = append(bannersByVer, b)
			}
		}

		var newestBanner model.Banner
		newestTime := time.Time{}

		for _, b := range bannersByVer {
			if b.UpdatedAt.After(newestTime) {
				newestBanner = b
				newestTime = b.UpdatedAt
			}
		}

		return &newestBanner, nil
	}
}

func UpdateBanner(bannerData model.BannerUpdateRequest, resp *response.Response) bool {

	var bannersByVer []model.Banner

	banners := pulls.GetBannersByFeatureId(bannerData.FeatureId)
	if len(banners) == 0 {
		resp.SetError(errs.GetErr(112))
		return false
	}

	for _, b := range banners {
		if reflect.DeepEqual(bannerData.TagIds, b.TagIds) {
			bannersByVer = append(bannersByVer, b)
		}
	}

	if len(bannersByVer) == 0 {
		resp.SetError(errs.GetErr(112))
		return false
	}

	var newestBanner model.Banner
	newestTime := time.Time{}
	oldestTime := bannersByVer[0].UpdatedAt
	oldestBanner := bannersByVer[0]

	if len(bannersByVer) > 1 {
		for _, b := range bannersByVer {
			if b.UpdatedAt.Sub(newestTime) > 0 {
				newestBanner = b
				newestTime = b.UpdatedAt
			}
		}
		for _, b := range bannersByVer {
			if b.UpdatedAt.Sub(oldestTime) < 0 {
				oldestBanner = b
				oldestTime = b.UpdatedAt
			}
		}

	} else if len(bannersByVer) == 1 {
		newestBanner = bannersByVer[0]
	}

	newestBanner.BannerItem = bannerData.BannerItem
	newestBanner.IsActive = bannerData.IsActive
	newestBanner.UpdatedAt = time.Now()

	id, err := postgres.InsertBanner(newestBanner)
	if err != nil {
		resp.SetError(errs.GetErr(99, err.Error()))
		return false
	}
	newestBanner.UpdatedAt = time.Now().UTC().Add(5 * time.Hour)
	newestBanner.BannerId = id
	pulls.AddPullBanner(newestBanner)

	if len(bannersByVer) >= 3 {

		err = postgres.DeleteBannerByID(oldestBanner.BannerId)
		if err != nil {
			resp.SetError(errs.GetErr(99, err.Error()))
			return false
		}

		pulls.DeleteBannerById(oldestBanner)
	}

	return true
}

func GetBannerVersions(bannerData model.BannerVersionsRequest) []model.Banner {

	var bannersByVer []model.Banner

	banners := pulls.GetBannersByFeatureId(bannerData.FeatureId)
	if len(banners) == 0 {
		return nil
	}

	for _, b := range banners {
		if reflect.DeepEqual(bannerData.TagIds, b.TagIds) {
			bannersByVer = append(bannersByVer, b)
		}
	}

	if len(bannersByVer) == 0 {
		return nil
	}

	return bannersByVer
}

func checkExistsBanner(data model.Banner) (int, bool) {

	banners := pulls.GetBannersByFeatureId(data.FeatureId)
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
