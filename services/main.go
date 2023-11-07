package services

import (
	"fmt"
	"recofiit/services/database"
)

func Register() {
	fmt.Println("Registering service")
	database.Init()
}
