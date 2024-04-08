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
		ctx.SetBodyString(service.Desc400)
		return
	}

	if ers := input.Validate(); ers != nil {
		log.Println("failed to validate create banner request")
		resp.SetErrors(ers)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(service.Desc400)
		return
	}

	id := service.CreateBanner(input, resp)
	if id == 0 {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(service.Desc500)
		return
	}

	var res = struct {
		BannerId int `json:"banner_id"`
	}{
		BannerId: id,
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBodyString(service.Desc201)

	resp.SetValue(res)
}

func GetBanner(resp *response.Response, ctx *fasthttp.RequestCtx) {
	var input model.BannerGetRequest
	var err error

	input.TagId, err = strconv.Atoi(string(ctx.QueryArgs().Peek("tag_id")))
	input.FeatureId, err = strconv.Atoi(string(ctx.QueryArgs().Peek("feature_id")))
	input.UseLastRevision, err = strconv.ParseBool(string(ctx.QueryArgs().Peek("use_last_revision")))

	if err != nil {
		log.Println("add banner error: " + err.Error())
		resp.SetError(errs.GetErr(100))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(service.Desc400)
	}

	if ers := input.Validate(); ers != nil {
		log.Println("failed to validate get banner request")
		resp.SetErrors(ers)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(service.Desc400)
		return
	}

	banner, err := service.GetBanner(input)
	if err != nil {
		resp.SetError(errs.GetErr(112))
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBodyString(service.Desc404)
		return
	}

	token := string(ctx.Request.Header.Peek("Authorization"))

	if !banner.IsActive && token != config.Cfg.Tokens.Admin {
		resp.SetError(errs.GetErr(112))
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBodyString(service.Desc404)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(service.Desc200)

	resp.SetValue(banner.BannerItem)
}

func UpdateBanner(resp *response.Response, ctx *fasthttp.RequestCtx) {
	var input model.BannerUpdateRequest

	bannerId, err := strconv.Atoi(string(ctx.QueryArgs().Peek("id")))

	if err = json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println("add banner error: " + err.Error())
		resp.SetError(errs.GetErr(100))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(service.Desc400)
		return
	}

	input.BannerId = bannerId

	if ers := input.Validate(); ers != nil {
		log.Println("failed to validate create banner request")
		resp.SetErrors(ers)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(service.Desc400)
		return
	}

	ok := service.UpdateBanner(input, resp)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(service.Desc500)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(service.Desc200)
}

func GetBannerVersions(resp *response.Response, ctx *fasthttp.RequestCtx) {
	var input model.BannerVersionsRequest

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println("add banner error: " + err.Error())
		resp.SetError(errs.GetErr(100))
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(service.Desc400)
		return
	}

	/*if ers := input.Validate(); ers != nil {
		log.Println("failed to validate create banner request")
		resp.SetErrors(ers)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString(service.Desc400)
		return
	}*/

	banners := service.GetBannerVersions(input)
	if len(banners) == 0 {
		resp.SetError(errs.GetErr(112))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(service.Desc500)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Версии баннеров пользователей")

	resp.SetValues(banners)
}
