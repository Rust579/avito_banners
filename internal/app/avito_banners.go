package app

import (
	"avito_banners/internal/config"
	"avito_banners/internal/handler"
	"avito_banners/internal/repo/postgres"
	"avito_banners/internal/server"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	if err := config.Init(); err != nil {
		log.Println("parsing configs error", err.Error())
		return
	}

	bts, _ := json.Marshal(config.Cfg)
	log.Println("________________________________CONFIG PARSING________________________________")
	log.Println(string(bts))

	if err := postgres.Init(); err != nil {
		log.Println("postgres init error", err.Error())
		return
	}

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
