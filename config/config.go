package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type BotConfig struct {
	BotToken                   string   `yaml:"bot_token"`
	BotCompanyName             string   `yaml:"bot_company_name"`
	DeliveryRatelimit          int      `yaml:"delivery_ratelimit"`
	SupportChatId              int64    `yaml:"support_chatid"`
	SupportLink                string   `yaml:"support_link"`
	DefaultConversationSla     int      `yaml:"default_sla"`
	StaleQuestionCheckInterval int      `yaml:"stale_question_check_interval"`
	StaleQuestionAfter         int      `yaml:"stale_question_after"`
	Agents                     []*Agent `yaml:"agents"`
	Db                         *Db      `yaml:"db"`
}

type Db struct {
	Host     string `yaml:"host"`
	DbName   string `yaml:"dbname"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SslMode  string `yaml:"sslmode"`
	Migrate  bool   `yaml:"migrate"`
}

type Agent struct {
	Name   string `yaml:"name"`
	ChatId int64  `yaml:"chatid"`
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
