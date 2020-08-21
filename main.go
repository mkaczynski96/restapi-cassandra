package main

import (
	"log"
	"restapi-cassandra/configs"
	"restapi-cassandra/controllers"
)

var cfg configs.Config

func main() {
	config := cfg.LoadConfig()
	connection, err := configs.BuildSession(config)
	if err != nil {
		log.Fatal(err)
	}
	controllers.RunApi(config, connection)
}
