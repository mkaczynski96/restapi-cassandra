package main

import (
	"acaisoft-mkaczynski-api/configs"
	"acaisoft-mkaczynski-api/controllers"
)

func main() {
	configs.SetupDBConnection()
	controllers.RunApi()
}
