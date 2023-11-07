package models

import (
	"gorm.io/gorm"
)

type sensorType string

const (
	FRONT_LIDAR            sensorType = "FRONT_LIDAR"
	FRONT_ULTRASONIC       sensorType = "FRONT_ULTRASONIC"
	REAR_ULTASONIC         sensorType = "REAR_ULTASONIC"
	WHEEL_SPEED            sensorType = "WHEEL_SPEED"
	GPS_LOCATION           sensorType = "GPS_LOCATION"
	GPS_SPEED              sensorType = "GPS_SPEED"
	GPS_DIRECTION          sensorType = "GPS_DIRECTION"
	MAGNETOMETER_DIRECTION sensorType = "MAGNETOMETER_DIRECTION"
)

type Car struct {
	gorm.Model
	Vin      string
	Name     string
	Sessions []Session `gorm:"many2many:car_sessions;"`
}

type Controller struct {
	gorm.Model
	Name                string `gorm:"type:varchar(255)"`
	Type                string `gorm:"type:varchar(255)"`
	Description         string `gorm:"type:text"`
	ControllerInstances []ControllerInstace
}

type Firmware struct {
	gorm.Model
	Version            string `gorm:"type:varchar(255)"`
	Description        string `gorm:"type:text"`
	ControllerInstaces []ControllerInstace
}

type Meassurement struct {
	CarSessionID int
	latency      int
	CreatedAt    string `gorm:"type:timestamptz"`
	CarSession   CarSession
}

type Sensor struct {
	gorm.Model
	ControllerInstaceID uint
	Name                string     `gorm:"type:varchar(255)"`
	SensorType          sensorType `gorm:"type:varchar(255)"`
	SensorData          []SensorData
	ControllerInstace   ControllerInstace
}

type Session struct {
	gorm.Model
	Name      string `gorm:"type:varchar(255)"`
	StartedAt string `gorm:"type:timestamptz"`
	EndedAt   string `gorm:"type:timestamptz"`
	Cars      []Car  `gorm:"many2many:car_sessions;"`
}

type SensorData struct {
	gorm.Model
	SensorID       uint
	MeassurementID uint
	Data1          string `gorm:"type:double precision"`
	Data2          string `gorm:"type:double precision"`
	Sensor         Sensor
}

type CarSession struct {
	gorm.Model
	CarID              uint
	SessionID          uint
	IsControlledByUser bool
	ControllerInstaces []ControllerInstace `gorm:"many2many:car_session_controllers;"`
	Meassurements      []Meassurement
	Session            Session
}

type CarSessionController struct {
	gorm.Model
	CarSessionID         uint
	ControllerInstanceID uint
}

type CarController struct {
	gorm.Model
	CarID                uint `gorm:"primaryKey"`
	ControllerInstanceID uint `gorm:"primaryKey"`
}

type ControllerInstace struct {
	gorm.Model
	FirmwareID   uint
	ControllerID uint
	name         string `gorm:"type:varchar(255)"`
	Firmware     Firmware
	Controller   Controller
}
