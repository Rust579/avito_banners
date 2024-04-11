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

	// Проверяем токен юзера или админа для эндпоинта /get-banner
	if !checktoken.CheckToken(ctx, resp) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Write(resp.FormResponse().Json())
		return
	}
	// Получение баннера
	if string(ctx.Path()) == "/get-banner" {
		GetBanner(resp, ctx)
		ctx.Write(resp.FormResponse().Json())
		return
	}

	// Проверяем токен админа для остальных эндпоинтов
	if !checktoken.CheckAdminToken(ctx, resp) {
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Write(resp.FormResponse().Json())
		return
	}

	switch string(ctx.Path()) {
	// Создание баннера
	case "/create-banner":
		CreateBanner(resp, ctx)
		// Обновление содержимого баннера с сохранением предыдущих двух версий
	case "/update-banner":
		UpdateBanner(resp, ctx)
		// Получение баннеров по фиче или тегу
	case "/get-banners":
		GetBanners(resp, ctx)
		// Удаление баннера по banner_id
	case "/delete-banner":
		DeleteBanner(resp, ctx)

		// Получение всех версий одного баннера
	case "/banner-versions":
		GetBannerVersions(resp, ctx)
		// Выбор версии баннера
	case "/set-banner-version":
		SetBannerVersion(resp, ctx)
		// Удаление баннеров по фиче или тегу
	case "/delete-banners":
		DeleteBanners(resp, ctx)

	default:
		log.Println("unknown request")
		resp.SetError(errs.GetErr(104))
	}

	ctx.Write(resp.FormResponse().Json())
}
