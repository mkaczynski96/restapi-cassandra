package main

import (
	"log"
	"restapi-cassandra/configs"
	"restapi-cassandra/controllers"
)

var cfg configs.Config

func main() {
	config := cfg.LoadConfig()
	err := configs.BuildSession(config, []string{config.Database.Address})
	if err != nil {
		log.Fatal(err)
	}
	controllers.RunApi(config)
}
