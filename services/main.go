package services

import (
	"fmt"
	"recofiit/services/communication"
	"recofiit/services/dataLogging"
	"recofiit/services/database"
	"recofiit/services/redis"
)

func Register() {
	fmt.Println("Registering services")
	database.Init()
	dataLogging.Init()

	communication.Init()
	redis.Init()
}
