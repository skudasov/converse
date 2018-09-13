package bot

import (
	"testing"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/f4hrenh9it/converse/db"
	"flag"
	"github.com/f4hrenh9it/converse/config"
	"os/exec"
	"os"
	"log"
	"time"
	"github.com/stretchr/testify/assert"
	"fmt"
)

var updChan chan tgbotapi.Update
var respChan chan tgbotapi.Chattable

func RunDb() {
	cmd := exec.Command("docker-compose", "-f", "docker-compose-pgonly.yaml", "up")
	cmd.Start()
}

func ClearData() {
	cmd := exec.Command("rm", "-rf", "pg_data")
	cmd.Start()
}

func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		log.Fatal(err)
	}
	RunDb()
	cfg := flag.String("config", "bot.yaml", "yaml configuration file")
	flag.Parse()

	if err := config.ParseBotConfig(*cfg); err != nil {
		log.Fatal(err)
	}
	tries := 10
	for i := 0; i < tries; i++ {
		time.Sleep(1 * time.Second)
		db.ConnectDb(config.B.Db)
		db.MigrateUp(config.B.Db, "file://migrations")
		db.RegisterAgents(config.B.Agents)
	}
	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updChan = make(chan tgbotapi.Update)

	NewStore(100200300)
	respChan = make(chan tgbotapi.Chattable)
	go StartReceiver(updChan, respChan, 10, 12)

	c := m.Run()
	ClearData()
	os.Exit(c)
}

func TestBot_NotInAnyConversation(t *testing.T) {
	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			MessageID: 10,
			Text:      "Hello",
			From: &tgbotapi.User{
				ID: 116664420,
			},
			Chat: &tgbotapi.Chat{
				ID: int64(116664420),
			},
		},
	}
	updChan <- upd
	res := <-respChan
	assert.Equal(t, res.(tgbotapi.MessageConfig), notInAnyConversation(116664420))
}

func TestBot_Start(t *testing.T) {
	upd := tgbotapi.Update{
		Message: &tgbotapi.Message{
			MessageID: 10,
			Text:      "/start",
			Entities: &[]tgbotapi.MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 6,
					User: &tgbotapi.User{
						ID: 116664420,
					},
				},
			},
			From: &tgbotapi.User{
				ID: 116664420,
			},
			Chat: &tgbotapi.Chat{
				ID: int64(116664420),
			},
		},
	}
	updChan <- upd
	res := <-respChan
	fmt.Printf("res: %s", res)
}
