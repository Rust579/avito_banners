package app

import (
	"avito_banners/internal/config"
	"avito_banners/internal/handler"
	"avito_banners/internal/repo/postgres"
	"avito_banners/internal/server"
	"avito_banners/internal/service/pulls"
	"encoding/json"
	"fmt"
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

	//TODO почистить память
	banners, err := postgres.GetAllBanners()
	if err != nil {
		log.Println("error get all banners from postgres", err.Error())
		return
	}

	pulls.InitBannersPulls(banners)

	srv := server.NewServer(config.Cfg.Service.Address)
	go func() {
		if err := srv.Run(handler.ServerHandler); err != nil {
			log.Printf(fmt.Sprintf("error occured while running http server: %s", err.Error()))
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
