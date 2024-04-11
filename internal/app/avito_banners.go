package app

import (
	"avito_banners/internal/config"
	"avito_banners/internal/handler"
	"avito_banners/internal/repo/postgres"
	"avito_banners/internal/server"
	"avito_banners/internal/service/pulls"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	if err := config.Init(); err != nil {
		log.Println("parsing configs error", err.Error())
		return
	}

	bts, _ := json.Marshal(config.Cfg)
	log.Println("________________________________CONFIG PARSING________________________________")
	log.Println(string(bts))

	var err error

	// Несколько попыток подключения к базе с задержкой в 5 секунд, на случай если база поднимется позже сервиса
	for i := 0; i <= 3; i++ {
		if err = postgres.Init(); err != nil {
			log.Println("postgres init error", err.Error())
			log.Println("waiting 5 second...")
		} else {
			break
		}
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Println("postgres init error", err.Error())
		return
	}

	// Достаем все баннеры с базы
	banners, err := postgres.GetAllBanners()
	if err != nil {
		log.Println("error get all banners from postgres", err.Error())
		return
	}

	// и собираем кэш со всеми баннерами
	pulls.InitBannersPulls(banners)

	// run http сервера
	srv := server.NewServer(config.Cfg.Service.Address)
	go func() {
		if err := srv.Run(handler.ServerHandler); err != nil {
			log.Println("error occurred while running http server: ", err.Error())
		}
	}()
	log.Printf("HTTP server started on addr  %s", config.Cfg.Service.Address)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(); err != nil {
		return
	}
}
