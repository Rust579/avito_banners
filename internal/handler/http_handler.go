package handler

import (
	"avito_banners/internal/errs"
	"avito_banners/internal/response"
	"github.com/valyala/fasthttp"
	"log"
)

func ServerHandler(ctx *fasthttp.RequestCtx) {

	ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, "*")
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodPost)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodGet)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodPatch)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowHeaders, fasthttp.HeaderContentType)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowHeaders, fasthttp.HeaderAuthorization)
	ctx.Response.Header.Set("Content-Type", "application/json")

	if ctx.IsOptions() {
		return
	}

	resp := response.InitResponse()

	switch string(ctx.Path()) {
	case "/xxx":
		AddBanner(resp, ctx)

	default:
		log.Println("unknown request")
		resp.SetError(errs.GetErr(104))
	}

	ctx.Write(resp.FormResponse().Json())
}
