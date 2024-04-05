package app

import (
	"avito_banners/internal/config"
	"avito_banners/internal/repo/postgres"
	"encoding/json"
	"log"
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
}
