package services

import (
	"fmt"
	"recofiit/services/communication"
	"recofiit/services/database"
	"recofiit/services/redis"
)

func Register() {
	fmt.Println("Registering services")
	database.Init()
	communication.Init()
	redis.Init()
}
