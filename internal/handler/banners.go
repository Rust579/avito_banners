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
	var input model.Banner

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println("add banner error: " + err.Error())
		resp.SetError(errs.GetErr(100))
		return
	}

	if !input.Validate() {
		log.Println("failed to validate create banner request")
		resp.SetError(errs.GetErr(400))
		return
	}

	id, err := service.AddBanner(input)
	if err != nil {
		log.Println("add banner error: " + err.Error())
		resp.SetError(errs.GetErr(99, err.Error()))
		return
	}

	var res = struct {
		Description string `json:"description"`
		BannerId    int    `json:"banner_id"`
	}{
		Description: "created",
		BannerId:    id,
	}

	resp.SetValue(res)
}

func GetBanner(resp *response.Response, ctx *fasthttp.RequestCtx) {

}
