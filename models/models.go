package models

import (
	"time"

	"gorm.io/gorm"
)

type SensorType string

const (
	GPS                 SensorType = "GPS"
	DISTANCE_ULTRASONIC SensorType = "DISTANCE_ULTRASONIC"
	DISTANCE_LIDAR      SensorType = "DISTANCE_LIDAR"
	SPEED_FRONT_LEFT    SensorType = "SPEED_FRONT_LEFT"
	SPEED_FRONT_RIGHT   SensorType = "SPEED_FRONT_RIGHT"
	SPEED_REAR_LEFT     SensorType = "SPEED_REAR_LEFT"
	SPEED_REAR_RIGHT    SensorType = "SPEED_REAR_RIGHT"
)

type Car struct {
	gorm.Model
	Vin      string
	Name     string
	Color    string
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
