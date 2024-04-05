package checktoken

import (
	"avito_banners/internal/config"
	"avito_banners/internal/errs"
	"avito_banners/internal/response"
	"github.com/valyala/fasthttp"
)

func CheckToken(ctx *fasthttp.RequestCtx, resp *response.Response) bool {
	token := string(ctx.Request.Header.Peek("Authorization"))

	if token == "" {
		resp.SetError(errs.GetErr(109))
		return false
	}

	if token != config.Cfg.Tokens.User && token != config.Cfg.Tokens.Admin {
		resp.SetError(errs.GetErr(108))
		return false
	}

	return true
}

func CheckAdminToken(ctx *fasthttp.RequestCtx, resp *response.Response) bool {
	token := string(ctx.Request.Header.Peek("Authorization"))

	if token == "" {
		resp.SetError(errs.GetErr(109))
		return false
	}

	if token != config.Cfg.Tokens.Admin {
		resp.SetError(errs.GetErr(110))
		return false
	}
	return true
}
