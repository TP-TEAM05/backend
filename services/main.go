package services

import (
	"fmt"
	"recofiit/services/communication"
	"recofiit/services/dataLogging"
	"recofiit/services/database"
	Logger "recofiit/services/logger"
	"recofiit/services/redis"
)

func Register() {
	fmt.Println("Registering services")
	Logger.Init()
	database.Init()
	dataLogging.Init()

	communication.Init()
	redis.Init()
}
