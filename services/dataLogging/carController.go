package dataLogging

import (
	"recofiit/models"
	"recofiit/services/database"
	"time"
)

type CarController struct{}

func NewCarController() *CarController {
	return &CarController{}
}

// CreateSessionAndCarSession creates a new session and car session, returning the car session ID
func (c *CarController) CreateSessionAndCarSession(carID uint, sessionName string) (uint, error) {
	// Create new session
	session := models.Session{
		Name:      sessionName,
		StartedAt: &time.Time{},
		EndedAt:   nil, // Set to nil for NULL in SQL
	}
	err := database.GetDB().Create(&session).Error
	if err != nil {
		return 0, err
	}

	// Create new car session
	carSession := models.CarSession{CarID: carID, SessionID: session.ID}
	err = database.GetDB().Create(&carSession).Error
	if err != nil {
		return 0, err
	}

	return carSession.ID, nil
}
