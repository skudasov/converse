package main

import (
	"flag"
	"github.com/f4hrenh9it/parley/bot"
	"github.com/f4hrenh9it/parley/config"
	"github.com/f4hrenh9it/parley/log"
	"github.com/f4hrenh9it/parley/db"
	_ "github.com/lib/pq"
)

func main() {

	cfg := flag.String("config", "bot.yaml", "yaml configuration file")
	flag.Parse()

	if err := config.ParseBotConfig(*cfg); err != nil {
		log.L.Fatal(err)
	}

	db.ConnectDb(config.B.Db)
	if config.B.Db.MigrateTestData {
		if err := db.MigrateUp(); err != nil {
			log.L.Fatalf("migration up failed: %s", err)
		}
	}

	bot.NewStore(config.B.SupportChatId)
	bot.StartNewBot(config.B.BotToken, config.B.DeliveryRatelimit, config.B.DefaultConversationSla)
}
