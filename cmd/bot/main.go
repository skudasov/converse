package main

import (
	"flag"
	"github.com/f4hrenh9it/converse/bot"
	"github.com/f4hrenh9it/converse/config"
	"github.com/f4hrenh9it/converse/log"
	"github.com/f4hrenh9it/converse/db"
	_ "github.com/lib/pq"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	cfg := flag.String("config", "bot.yaml", "yaml configuration file")
	flag.Parse()

	if err := config.ParseBotConfig(*cfg); err != nil {
		log.L.Fatal(err)
	}

	db.ConnectDb(config.B.Db)
	if config.B.Db.Migrate {
		if err := db.MigrateUp(config.B.Db, "file:///migrations"); err != nil {
			log.L.Fatalf("migration up failed: %s", err)
		}
	}
	if err := db.RegisterAgents(config.B.Agents); err != nil {
		log.L.Errorf(db.AgentRegisterErr, err)
	}

	bot.NewStore(config.B.SupportChatId)

	log.L.Debugf("using token: %s", config.B.BotToken)
	b, err := tgbotapi.NewBotAPI(config.B.BotToken)
	if err != nil {
		log.L.Fatal(err)
	}
	b.Debug = false
	log.L.Debugf("Authorized on account %s", b.Self.UserName)

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60

	updChan, err := b.GetUpdatesChan(ucfg)
	if err != nil {
		log.L.Fatalf("failed to get bot updates channel: %s", err)
	}

	msgChan := make(chan tgbotapi.Chattable)

	go bot.StartSender(b, msgChan)
	bot.StartReceiver(updChan, msgChan, config.B.DeliveryRatelimit, config.B.DefaultConversationSla)
}
