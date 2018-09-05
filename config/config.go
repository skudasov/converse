package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Db struct {
	Host     string `yaml:"host"`
	DbName   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SslMode  string `yaml:"sslmode"`
	MigrateTestData            bool   `yaml:"migrate_test_data"`
}

type BotConfig struct {
	BotToken                   string `yaml:"bot_token"`
	DeliveryRatelimit          int    `yaml:"delivery_ratelimit"`
	SupportChatId              int64  `yaml:"support_chatid"`
	DefaultConversationSla     int    `yaml:"default_sla"`
	StaleQuestionCheckInterval int    `yaml:"stale_question_check_interval"`
	StaleQuestionAfter         int    `yaml:"stale_question_after"`
	Db                         *Db    `yaml:"db"`
}

var B *BotConfig

func ParseBotConfig(f string) error {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	B = &BotConfig{}
	err = yaml.Unmarshal(b, B)
	if err != nil {
		return err
	}
	return nil
}
