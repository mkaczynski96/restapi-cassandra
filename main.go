package main

import (
	"log"
	"restapi-cassandra/configs"
	"restapi-cassandra/controllers"
)

func main() {
	config := configs.LoadConfig()
	connection, err := configs.BuildSession(config)
	if err != nil {
		log.Fatal(err)
	}
	controllers.NewApi(config, connection)
}
