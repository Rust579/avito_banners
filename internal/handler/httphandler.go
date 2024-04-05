package handler

import (
	"avito_banners/internal/errs"
	"avito_banners/internal/pkg/checktoken"
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

	if err := checktoken.CheckToken(ctx); err != nil {
		resp.SetError(errs.GetErr(401))
		ctx.Write(resp.FormResponse().Json())
		return
	}

	if string(ctx.Path()) == "/get-banner" {
		GetBanner(resp, ctx)
	}

	if err := checktoken.CheckAdminToken(ctx); err != nil {
		resp.SetError(errs.GetErr(403))
		ctx.Write(resp.FormResponse().Json())
		return
	}

	switch string(ctx.Path()) {
	case "/create-banner":
		CreateBanner(resp, ctx)

	default:
		log.Println("unknown request")
		resp.SetError(errs.GetErr(104))
	}

	ctx.Write(resp.FormResponse().Json())
}
