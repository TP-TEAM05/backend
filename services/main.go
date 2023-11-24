package services

import (
	"fmt"
	"recofiit/services/communication"
	"recofiit/services/database"
)

func Register() {
	fmt.Println("Registering services")
	database.Init()
	communication.Init()
}
