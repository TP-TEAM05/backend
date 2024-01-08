package models

import (
	"time"

	"gorm.io/gorm"
)

type SensorType string

const (
	FRONT_LIDAR            SensorType = "FRONT_LIDAR"
	FRONT_ULTRASONIC       SensorType = "FRONT_ULTRASONIC"
	REAR_ULTRASONIC        SensorType = "REAR_ULTRASONIC"
	WHEEL_SPEED            SensorType = "WHEEL_SPEED"
	GPS_LOCATION           SensorType = "GPS_LOCATION"
	GPS_SPEED              SensorType = "GPS_SPEED"
	GPS_DIRECTION          SensorType = "GPS_DIRECTION"
	MAGNETOMETER_DIRECTION SensorType = "MAGNETOMETER_DIRECTION"
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
	ControllerInstances []ControllerInstance
}

type Firmware struct {
	gorm.Model
	Version             string `gorm:"type:varchar(255)"`
	Description         string `gorm:"type:text"`
	ControllerInstances []ControllerInstance
}

type Measurement struct {
	CarSessionID uint
	SensorID     uint
	latency      int
	CreatedAt    *time.Time `gorm:"type:timestamptz"`
	Data1        float32    `gorm:"type:double precision"`
	Data2        float32    `gorm:"type:double precision"`
	Sensor       Sensor
	CarSession   CarSession
}

type Sensor struct {
	gorm.Model
	ControllerInstanceID uint
	Name                 string     `gorm:"type:varchar(255)"`
	SensorType           SensorType `gorm:"type:varchar(255)"`
	ControllerInstance   ControllerInstance
}

type Session struct {
	gorm.Model
	Name      string     `gorm:"type:varchar(255)"`
	StartedAt *time.Time `gorm:"type:timestamptz"`
	EndedAt   *time.Time `gorm:"type:timestamptz"`
	Cars      []Car      `gorm:"many2many:car_sessions;"`
}

type CarSession struct {
	gorm.Model
	CarID               uint
	SessionID           uint
	IsControlledByUser  bool
	ControllerInstances []ControllerInstance `gorm:"many2many:car_session_controllers;"`
	Meassurements       []Measurement
	Session             Session
}

type CarSessionController struct {
	gorm.Model
	CarSessionID         uint
	ControllerInstanceID uint
}

type CarController struct {
	gorm.Model
	CarID                uint
	ControllerInstanceID uint
}

type ControllerInstance struct {
	gorm.Model
	FirmwareID   uint
	ControllerID uint
	name         string `gorm:"type:varchar(255)"`
	Firmware     Firmware
	Controller   Controller
}
