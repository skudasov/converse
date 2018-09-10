package main

import (
	"flag"
	"github.com/f4hrenh9it/converse/bot"
	"github.com/f4hrenh9it/converse/config"
	"github.com/f4hrenh9it/converse/log"
	"github.com/f4hrenh9it/converse/db"
	_ "github.com/lib/pq"
)

func main() {

	cfg := flag.String("config", "bot.yaml", "yaml configuration file")
	flag.Parse()

	if err := config.ParseBotConfig(*cfg); err != nil {
		log.L.Fatal(err)
	}

	db.ConnectDb(config.B.Db)
	if config.B.Db.Migrate {
		if err := db.MigrateUp(config.B.Db); err != nil {
			log.L.Fatalf("migration up failed: %s", err)
		}
	}
	if err := db.RegisterAgents(config.B.Agents); err != nil {
		log.L.Errorf(db.AgentRegisterErr, err)
	}

	bot.NewStore(config.B.SupportChatId)
	bot.StartNewBot(config.B.BotToken, config.B.DeliveryRatelimit, config.B.DefaultConversationSla)
}
