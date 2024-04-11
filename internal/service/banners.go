package service

import (
	"avito_banners/internal/errs"
	"avito_banners/internal/model"
	"avito_banners/internal/repo/postgres"
	"avito_banners/internal/response"
	"avito_banners/internal/service/pulls"
	"errors"
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

	// Если такой баннер уже существует - ошибка
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

	// Если нужна актуальная информация по баннеру - идем сразу в базу
	if bannerData.UseLastRevision {
		banner, err := postgres.FindBannerByFeatureAndTagId(bannerData)
		if err != nil {
			return nil, err
		}

		return banner, nil

	} else {
		// Иначе достаем баннер из кэша
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

		// Выбираем последнюю версию
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

	// Далее определяем из какой версии баннера создать новую и какую версию можно удалить:
	var newestBanner model.Banner
	newestTime := time.Time{}
	oldestTime := bannersByVer[0].UpdatedAt
	oldestBanner := bannersByVer[0]

	// Если версий несколько:
	if len(bannersByVer) > 1 {
		// Определяем последнюю версию
		for _, b := range bannersByVer {
			if b.UpdatedAt.Sub(newestTime) > 0 {
				newestBanner = b
				newestTime = b.UpdatedAt
			}
		}
		// Определяем самую старую версию
		for _, b := range bannersByVer {
			if b.UpdatedAt.Sub(oldestTime) < 0 {
				oldestBanner = b
				oldestTime = b.UpdatedAt
			}
		}

		// Если версия одна, просто создаем по ней новую версию
	} else if len(bannersByVer) == 1 {
		newestBanner = bannersByVer[0]
	}

	newestBanner.BannerItem = bannerData.BannerItem
	newestBanner.IsActive = bannerData.IsActive
	newestBanner.UpdatedAt = time.Now().UTC()

	// Добавляем новую версию
	id, err := postgres.InsertBanner(newestBanner)
	if err != nil {
		resp.SetError(errs.GetErr(99, err.Error()))
		return false
	}

	newestBanner.BannerId = id
	pulls.AddPullBanner(newestBanner)

	// Если версий 3 или более - удаляем самую старую версию, чтобы версий было всегда не более 3
	if len(bannersByVer) >= 3 {

		if err = postgres.DeleteBannerByID(oldestBanner.BannerId); err != nil {
			resp.SetError(errs.GetErr(99, err.Error()))
			return false
		}

		pulls.DeleteBannerById(oldestBanner)
	}

	return true
}

// GetBannerVersions Возвращает все версии баннера
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

// SetBannerVersion Выбор какой-либо версии баннера как основной
func SetBannerVersion(input model.BannerIdRequest, resp *response.Response) error {

	banner := pulls.GetBannerById(input.BannerId)
	if banner == nil {
		resp.SetError(errs.GetErr(112))
		return errors.New("banner not found")
	}

	banner.UpdatedAt = time.Now().UTC()

	if err := postgres.SetNewBannerVersionByID(banner.BannerId, banner.UpdatedAt); err != nil {
		resp.SetError(errs.GetErr(99, err.Error()))
		return err
	}
	pulls.UpdatePullBanner(*banner)

	return nil
}

func GetBanners(input model.BannersGetRequest) ([]model.Banner, int, error) {

	banners, count, err := postgres.FindBannersByParams(input.TagId, input.FeatureId, input.Offset, input.Limit)
	if err != nil {
		return nil, count, err
	}

	return banners, count, nil
}

func DeleteBanner(input model.BannerIdRequest) error {

	if err := postgres.DeleteBannerByID(input.BannerId); err != nil {
		return err
	}

	banner := pulls.GetBannerById(input.BannerId)
	if banner == nil {
		return errors.New("banner not found")
	}
	pulls.DeleteBannerById(*banner)

	return nil
}

func DeleteBanners(input model.BannersDeleteRequest, resp *response.Response) ([]int, error) {

	var deletedBanners []model.Banner
	var deletedBannerIds []int
	var err error

	if input.TagId != 0 {
		deletedBanners, err = postgres.DeleteBannersByTagId(input.TagId)
		if err != nil {
			resp.SetError(errs.GetErr(99, err.Error()))
			return nil, err
		}
	}

	if input.FeatureId != 0 {
		deletedBanners, err = postgres.DeleteBannersByFeatureId(input.FeatureId)
		if err != nil {
			resp.SetError(errs.GetErr(99, err.Error()))
			return nil, err
		}
	}

	if len(deletedBanners) == 0 {
		resp.SetError(errs.GetErr(115))
		return nil, errors.New("banners not found")
	}

	for _, b := range deletedBanners {
		pulls.DeleteBannerById(b)
	}

	for _, b := range deletedBanners {
		deletedBannerIds = append(deletedBannerIds, b.BannerId)
	}

	return deletedBannerIds, nil
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
