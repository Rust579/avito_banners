package handler

import (
	"avito_banners/internal/errs"
	"avito_banners/internal/pkg/checktoken"
	"avito_banners/internal/response"
	"avito_banners/internal/service"
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

	// Проверяем токен юзера или админа для эндпоинта /get-banner
	if !checktoken.CheckToken(ctx, resp) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetBodyString(service.Desc401)
		ctx.Write(resp.FormResponse().Json())
		return
	}

	if string(ctx.Path()) == "/user_banner" {
		GetBanner(resp, ctx)
		ctx.Write(resp.FormResponse().Json())
		return
	}

	// Проверяем токен админа для остальных эндпоинтов
	if !checktoken.CheckAdminToken(ctx, resp) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.SetBodyString(service.Desc403)
		ctx.Write(resp.FormResponse().Json())
		return
	}

	switch string(ctx.Path()) {
	case "/create-banner":
		CreateBanner(resp, ctx)
	case "/banner":
		UpdateBanner(resp, ctx)
	case "/banners":
		GetBanners(resp, ctx)
	case "/delete":
		DeleteBanner(resp, ctx)

	case "/banner-versions":
		GetBannerVersions(resp, ctx)
	case "/set-banner-version":
		SetBannerVersion(resp, ctx)

	default:
		log.Println("unknown request")
		resp.SetError(errs.GetErr(104))
	}

	ctx.Write(resp.FormResponse().Json())
}
