package handler

import (
	"avito_banners/internal/config"
	"avito_banners/internal/errs"
	"avito_banners/internal/model"
	"avito_banners/internal/response"
	"avito_banners/internal/service"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
)

func CreateBanner(resp *response.Response, ctx *fasthttp.RequestCtx) {
	var input model.BannerAddRequest

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println("add banner error: " + err.Error())
		resp.SetError(errs.GetErr(100))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if ers := input.Validate(); ers != nil {
		resp.SetErrors(ers)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	id := service.CreateBanner(input, resp)
	if id == 0 {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	var res = struct {
		BannerId int `json:"banner_id"`
	}{
		BannerId: id,
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)

	resp.SetValue(res)
}

func GetBanner(resp *response.Response, ctx *fasthttp.RequestCtx) {
	var input model.BannerGetRequest
	var err error

	input.TagId, err = strconv.Atoi(string(ctx.QueryArgs().Peek("tag_id")))
	input.FeatureId, err = strconv.Atoi(string(ctx.QueryArgs().Peek("feature_id")))
	input.UseLastRevision, err = strconv.ParseBool(string(ctx.QueryArgs().Peek("use_last_revision")))

	if err != nil {
		log.Println("get banner error: " + err.Error())
		resp.SetError(errs.GetErr(100))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if ers := input.Validate(); ers != nil {
		resp.SetErrors(ers)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	banner, err := service.GetBanner(input)
	if err != nil {
		resp.SetError(errs.GetErr(112))
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	token := string(ctx.Request.Header.Peek("Authorization"))

	if !banner.IsActive && token != config.Cfg.Tokens.Admin {
		resp.SetError(errs.GetErr(112))
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)

	resp.SetValue(banner.BannerItem)
}

func UpdateBanner(resp *response.Response, ctx *fasthttp.RequestCtx) {
	var input model.BannerUpdateRequest

	bannerId, err := strconv.Atoi(string(ctx.QueryArgs().Peek("id")))

	if err = json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println("update banner error: " + err.Error())
		resp.SetError(errs.GetErr(100))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	input.BannerId = bannerId

	if ers := input.Validate(); ers != nil {
		resp.SetErrors(ers)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	ok := service.UpdateBanner(input, resp)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func GetBannerVersions(resp *response.Response, ctx *fasthttp.RequestCtx) {
	var input model.BannerVersionsRequest

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println("get banner versions error: " + err.Error())
		resp.SetError(errs.GetErr(100))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if ers := input.Validate(); ers != nil {
		resp.SetErrors(ers)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	banners := service.GetBannerVersions(input)
	if len(banners) == 0 {
		resp.SetError(errs.GetErr(112))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)

	resp.SetValues(banners)
}

func SetBannerVersion(resp *response.Response, ctx *fasthttp.RequestCtx) {
	var input model.BannerIdRequest

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println("set banner version error: " + err.Error())
		resp.SetError(errs.GetErr(100))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if ers := input.Validate(); ers != nil {
		resp.SetErrors(ers)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	err := service.SetBannerVersion(input, resp)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func GetBanners(resp *response.Response, ctx *fasthttp.RequestCtx) {
	var input model.BannersGetRequest
	var err error

	input.TagId, _ = strconv.Atoi(string(ctx.QueryArgs().Peek("tag_id")))
	input.FeatureId, _ = strconv.Atoi(string(ctx.QueryArgs().Peek("feature_id")))
	input.Limit, _ = strconv.Atoi(string(ctx.QueryArgs().Peek("limit")))
	input.Offset, _ = strconv.Atoi(string(ctx.QueryArgs().Peek("offset")))

	if ers := input.Validate(); ers != nil {
		resp.SetErrors(ers)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	banners, count, err := service.GetBanners(input)
	if err != nil {
		resp.SetError(errs.GetErr(99, err.Error()))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	if len(banners) == 0 {
		resp.SetError(errs.GetErr(115))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	var res = struct {
		Banners []*model.Banner `json:"banners"`
		Count   int             `json:"count"`
	}{
		Banners: banners,
		Count:   count,
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	resp.SetValue(res)
}

func DeleteBanner(resp *response.Response, ctx *fasthttp.RequestCtx) {
	var input model.BannerIdRequest
	var err error

	input.BannerId, err = strconv.Atoi(string(ctx.QueryArgs().Peek("id")))
	if err != nil {
		resp.SetError(errs.GetErr(100))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if ers := input.Validate(); ers != nil {
		resp.SetErrors(ers)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	err = service.DeleteBanner(input)
	if err != nil {
		resp.SetError(errs.GetErr(99, err.Error()))
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
	//TODO почему-то не пишется респонс пустой
}

func DeleteBanners(resp *response.Response, ctx *fasthttp.RequestCtx) {
	var input model.BannersDeleteRequest
	var err error

	input.TagId, err = strconv.Atoi(string(ctx.QueryArgs().Peek("tag_id")))
	input.FeatureId, err = strconv.Atoi(string(ctx.QueryArgs().Peek("feature_id")))

	if ers := input.Validate(); ers != nil {
		resp.SetErrors(ers)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	err = service.DeleteBanners(input)
	if err != nil {
		resp.SetError(errs.GetErr(112))
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
}
