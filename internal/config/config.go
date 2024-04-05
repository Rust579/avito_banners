package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Configs struct {
	Service  Service  `json:"service"`
	Postgres Postgres `json:"postgres"`
	Tokens   Tokens   `json:"tokens"`
}

type Service struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Postgres struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

type Tokens struct {
	Admin string `json:"admin"`
	User  string `json:"user"`
}

var Cfg Configs

func Init() error {

	if !fileExists("env.json") {
		return errors.New("file 'env.json' not found")
	}

	envFile, err := os.Open("env.json")
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(envFile)

	if er := decoder.Decode(&Cfg); er != nil {
		return er
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
