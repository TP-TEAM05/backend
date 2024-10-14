package dataLogging

import (
	"recofiit/models"
	"recofiit/services/database"
	"time"

	"github.com/getsentry/sentry-go"
)

type CarController struct{}

func NewCarController() *CarController {
	return &CarController{}
}

func (c *CarController) CreateSessionAndCarSession(carID uint, sessionName string) (uint, error) {
	var session models.Session
	var carSession models.CarSession

	if err := database.GetDB().Where("car_id = ?", carID).First(&carSession).Error; err != nil {
		session, err = c.CreateSession(carID, "Test session")
		if err != nil {
			sentry.CaptureException(err)
			return 0, err
		}
		carSession, err = c.CreateCarSession(carID, session.ID)
		if err != nil {
			sentry.CaptureException(err)
			return 0, err
		}
	}
	return carSession.ID, nil
}

func (c *CarController) CreateCarSession(carID uint, sessionId uint) (models.CarSession, error) {
	carSession := models.CarSession{CarID: carID, SessionID: sessionId}
	var err = database.GetDB().Create(&carSession).Error
	if err != nil {
		sentry.CaptureException(err)
		return carSession, err
	}

	return carSession, nil
}

func (c *CarController) CreateSession(carID uint, sessionName string) (models.Session, error) {
	session := models.Session{
		Name:      sessionName,
		StartedAt: &time.Time{},
		EndedAt:   nil,
	}
	err := database.GetDB().Create(&session).Error
	if err != nil {
		sentry.CaptureException(err)
		return session, err
	}

	return session, nil
}
