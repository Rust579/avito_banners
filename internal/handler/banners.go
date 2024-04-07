package handler

import (
	"avito_banners/internal/errs"
	"avito_banners/internal/model"
	"avito_banners/internal/response"
	"avito_banners/internal/service"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
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

	id := service.AddBanner(input, resp)
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

}
