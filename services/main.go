package services

import (
	"fmt"
	"recofiit/services/database"
	"recofiit/services/redis"
)

func Register() {
	fmt.Println("Registering service")
	database.Init()
	redis.Init()
}
