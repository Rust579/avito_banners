package checktoken

import (
	"avito_banners/internal/config"
	"errors"
	"github.com/valyala/fasthttp"
)

func CheckToken(ctx *fasthttp.RequestCtx) error {
	token := string(ctx.Request.Header.Peek("Authorization"))

	if token != config.Cfg.Tokens.User && token != config.Cfg.Tokens.Admin {
		return errors.New("пользователь не авторизован")
	}

	return nil
}

func CheckAdminToken(ctx *fasthttp.RequestCtx) error {
	token := string(ctx.Request.Header.Peek("Authorization"))

	if token != config.Cfg.Tokens.Admin {
		return errors.New("пользователь не авторизован")
	}
	return nil
}
