package configs

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Api struct {
		Port string `yaml:"port" ,envconfig:"API_PORT"`
	} `yaml:"api"`
	Database struct {
		Address      string `yaml:"address" ,envconfig:"DB_ADDRESS"`
		Port         int    `yaml:"port" ,envconfig:"DB_PORT"`
		ProtoVersion int    `yaml:"protoVersion" ,envconfig:"DB_PROTOVERSION"`
		Keyspace     string `yaml:"keyspace" ,envconfig:"DB_KEYSPACE"`
		TableName    string `yaml:"tableName" ,envconfig:"DB_TABLENAME"`
	} `yaml:"database"`
	Mail struct {
		MessageExpirationSeconds int    `yaml:"messageExpirationSeconds" ,envconfig:"MAIL_MESSAGEEXPIRATIONSECONDS"`
		Username                 string `yaml:"username" ,envconfig:"MAIL_USERNAME"`
		Password                 string `yaml:"password" ,envconfig:"MAIL_PASSWORD"`
		Host                     string `yaml:"host" ,envconfig:"MAIL_HOST"`
		Port                     int    `yaml:"port" ,envconfig:"MAIL_PORT"`
	} `yaml:"mail"`
}

func LoadConfig() Config {
	cfg := Config{}
	readFile(&cfg)
	readEnv(&cfg)
	log.Printf("Config: %+v", cfg)
	return cfg
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readFile(cfg *Config) {
	f, err := os.Open("./configs/config.yml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
