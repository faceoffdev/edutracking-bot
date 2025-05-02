package main

import (
	"github.com/faceoffdev/edutracking-bot/internal/bot"
	"github.com/faceoffdev/edutracking-bot/internal/config"
	"github.com/faceoffdev/edutracking-bot/internal/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	err = db.Init(cfg.Postgres)
	if err != nil {
		panic(err)
	}

	if err = bot.New().Run(cfg); err != nil {
		panic(err)
	}
}
