package tests

import (
	"avito_banners/internal/model"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"strconv"
	"testing"
)

var bodies []model.BannerAddRequest
var bannerIds []int

// Тесты всех эндпоинтов >>>>>
// База должна быть пустая

// create-banner >>>>>
func TestCreateBanner(t *testing.T) {

	uri := "/create-banner"
	headers := map[string]string{
		"Authorization": "tokenforadmin",
	}

	var res []struct {
		BannerId int `json:"banner_id"`
	}

	for i := 1; i <= 3; i++ {
		body := model.BannerAddRequest{
			TagIds:    []int{i*2 - 1, i * 2},
			FeatureId: i,
			IsActive:  true,
			BannerItem: map[string]interface{}{
				"text":  strconv.Itoa(i),
				"title": strconv.Itoa(i),
				"url":   strconv.Itoa(i),
			},
		}
		resp, err := SendHTTPRequest(uri, fasthttp.MethodPost, headers, fasthttp.StatusCreated, body)
		if err != nil {
			t.Fatalf("send request error: %v", err)
		}
		if len(resp.Errors) != 0 {
			t.Fatalf("response contain errors: %v", resp.Errors)
		}
		data, err := json.Marshal(resp.Values)
		if err != nil {
			t.Fatalf("marshal response.values error: %v", err)
		}
		err = json.Unmarshal(data, &res)
		if err != nil {
			t.Fatalf("unmarshal response.values error: %v", err)
		}
		for _, item := range res {
			bannerIds = append(bannerIds, item.BannerId)
		}

		bodies = append(bodies, body)
	}

	fmt.Println("Успешное создание баннеров:", bannerIds)
}

// get-banner >>>>>
func TestGetBanner(t *testing.T) {

	var res []struct {
		Text  string `json:"text"`
		Title string `json:"title"`
		Url   string `json:"url"`
	}
	var results []interface{}

	for _, b := range bodies {

		uri := "/get-banner?tag_id=" + strconv.Itoa(b.TagIds[0]) + "&feature_id=" + strconv.Itoa(b.FeatureId) + "&use_last_revision=true"
		headers := map[string]string{
			"Authorization": "tokenforadmin",
		}

		resp, err := SendHTTPRequest(uri, fasthttp.MethodGet, headers, fasthttp.StatusOK, nil)
		if err != nil {
			t.Fatalf("send request error: %v", err)
		}
		if len(resp.Errors) != 0 {
			t.Fatalf("response contain errors: %v", resp.Errors)
		}
		data, err := json.Marshal(resp.Values)
		if err != nil {
			t.Fatalf("marshal response.values error: %v", err)
		}
		err = json.Unmarshal(data, &res)
		if err != nil {
			t.Fatalf("unmarshal response.values error: %v", err)
		}

		if (res[0].Url != b.BannerItem["url"]) || (res[0].Text != b.BannerItem["text"]) || (res[0].Title != b.BannerItem["title"]) {
			t.Fatal("содержимое баннера не соответствует начальным данным")
		}

		results = append(results, res[0])
	}

	fmt.Println("Успешное получение баннеров:", results)
}

// update-banner >>>>>
func TestUpdateBanner(t *testing.T) {

	// На обновление отдаем первый элемент ранее созданного bodies
	body := model.BannerUpdateRequest{
		TagIds:    bodies[0].TagIds,
		FeatureId: bodies[0].FeatureId,
		IsActive:  true,
		BannerItem: map[string]interface{}{
			"text":  "updates_text",
			"title": "updates_title",
			"url":   "updates_url",
		},
	}

	uri := "/update-banner?id=" + strconv.Itoa(bannerIds[0])
	headers := map[string]string{
		"Authorization": "tokenforadmin",
	}

	resp, err := SendHTTPRequest(uri, fasthttp.MethodPatch, headers, fasthttp.StatusOK, body)
	if err != nil {
		t.Fatalf("send request error: %v", err)
	}
	if len(resp.Errors) != 0 {
		t.Fatalf("response contain errors: %v", resp.Errors)
	}

	fmt.Println("Успешное обновление баннера")
}

// get-banners >>>>>
func TestGetBanners(t *testing.T) {

	var res []struct {
		Banners []model.Banner `json:"banners"`
		Count   int            `json:"count"`
	}

	uri := "/get-banners?tag_id=" + strconv.Itoa(bodies[0].TagIds[0]) + "&feature_id=" + strconv.Itoa(bodies[0].FeatureId) + "&limit=1" + "&offset=0"
	headers := map[string]string{
		"Authorization": "tokenforadmin",
	}

	resp, err := SendHTTPRequest(uri, fasthttp.MethodGet, headers, fasthttp.StatusOK, nil)
	if err != nil {
		t.Fatalf("send request error: %v", err)
	}
	if len(resp.Errors) != 0 {
		t.Fatalf("response contain errors: %v", resp.Errors)
	}
	data, err := json.Marshal(resp.Values)
	if err != nil {
		t.Fatalf("marshal response.values error: %v", err)
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		t.Fatalf("unmarshal response.values error: %v", err)
	}
	if (res[0].Banners[0].FeatureId != bodies[0].FeatureId) || (res[0].Banners[0].TagIds[0] != bodies[0].TagIds[0]) || res[0].Count != 1 {
		t.Fatal("получен не корректный баннер")
	}

	fmt.Println("Успешное получение баннеров:", res)
}

// get-banner-versions >>>>>
func TestGetBannerVersions(t *testing.T) {

	var res []model.Banner

	body := model.BannerVersionsRequest{
		TagIds:    bodies[0].TagIds,
		FeatureId: bodies[0].FeatureId,
	}

	uri := "/get-banner-versions"
	headers := map[string]string{
		"Authorization": "tokenforadmin",
	}

	resp, err := SendHTTPRequest(uri, fasthttp.MethodGet, headers, fasthttp.StatusOK, body)
	if err != nil {
		t.Fatalf("send request error: %v", err)
	}
	if len(resp.Errors) != 0 {
		t.Fatalf("response contain errors: %v", resp.Errors)
	}
	data, err := json.Marshal(resp.Values)
	if err != nil {
		t.Fatalf("marshal response.values error: %v", err)
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		t.Fatalf("unmarshal response.values error: %v", err)
	}
	if len(res) != 2 {
		t.Fatal("получено неверное количество баннеров")
	}

	fmt.Println("Успешное получение версий баннеров:", res)
}

// set-banner-version >>>>>
func TestSetBannerVersion(t *testing.T) {

	body := model.BannerIdRequest{
		BannerId: bannerIds[0],
	}

	uri := "/set-banner-version"
	headers := map[string]string{
		"Authorization": "tokenforadmin",
	}

	resp, err := SendHTTPRequest(uri, fasthttp.MethodPatch, headers, fasthttp.StatusOK, body)
	if err != nil {
		t.Fatalf("send request error: %v", err)
	}
	if len(resp.Errors) != 0 {
		t.Fatalf("response contain errors: %v", resp.Errors)
	}

	fmt.Println("Успешный выбор версии баннера")
}

// delete-banner >>>>>
func TestDeleteBanner(t *testing.T) {

	uri := "/delete-banner?id=" + strconv.Itoa(bannerIds[0])
	headers := map[string]string{
		"Authorization": "tokenforadmin",
	}

	resp, err := SendHTTPRequest(uri, fasthttp.MethodDelete, headers, fasthttp.StatusNoContent, nil)
	if err != nil {
		t.Fatalf("send request error: %v", err)
	}
	if len(resp.Errors) != 0 {
		t.Fatalf("response contain errors: %v", resp.Errors)
	}

	fmt.Println("Успешное удаление баннера")
}

// delete-banners >>>>>
func TestDeleteBanners(t *testing.T) {

	for i := 1; i <= 2; i++ {
		uri := "/delete-banners?feature_id=" + strconv.Itoa(i)
		headers := map[string]string{
			"Authorization": "tokenforadmin",
		}

		resp, err := SendHTTPRequest(uri, fasthttp.MethodDelete, headers, fasthttp.StatusOK, nil)
		if err != nil {
			t.Fatalf("send request error: %v", err)
		}
		if len(resp.Errors) != 0 {
			t.Fatalf("response contain errors: %v", resp.Errors)
		}
	}

	uri := "/delete-banners?tag_id=5"
	headers := map[string]string{
		"Authorization": "tokenforadmin",
	}

	resp, err := SendHTTPRequest(uri, fasthttp.MethodDelete, headers, fasthttp.StatusOK, nil)
	if err != nil {
		t.Fatalf("send request error: %v", err)
	}
	if len(resp.Errors) != 0 {
		t.Fatalf("response contain errors: %v", resp.Errors)
	}

	fmt.Println("Успешное удаление баннеров")
}
