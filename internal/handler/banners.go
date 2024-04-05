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

func AddBanner(resp *response.Response, ctx *fasthttp.RequestCtx) {
	var input model.Banner

	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
		log.Println("add banner error: " + err.Error())
		resp.SetError(errs.GetErr(100))
		return
	}

	if ers := input.Validate(); ers != nil {
		log.Println("failed to validate add banner request", ers)
		resp.SetErrors(ers)
		return
	}

	if err := service.AddBanner(input); err != nil {
		log.Println("add banner error: " + err.Error())
		resp.SetError(errs.GetErr(99, err.Error()))
		return
	}

	resp.SetValue("all right")
}
