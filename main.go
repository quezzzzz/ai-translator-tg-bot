package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"tg_bot/config"
	"tg_bot/internal/tg_bot"
)

func main() {
	cfg, err := config.ParseConfig("config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	b, err := tg_bot.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	b.StartBot()

}
