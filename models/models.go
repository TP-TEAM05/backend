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

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Car struct {
	BaseModel
	Vin      string    `json:"vin"`
	Name     string    `json:"name"`
	Color    string    `json:"color"`
	Sessions []Session `gorm:"many2many:car_sessions;" json:"sessions"`
}

type Controller struct {
	BaseModel
	Name                string               `gorm:"type:varchar(255)" json:"name"`
	Type                string               `gorm:"type:varchar(255)" json:"type"`
	Description         string               `gorm:"type:text" json:"description"`
	ControllerInstances []ControllerInstance `json:"controller_instances"`
}

type Firmware struct {
	BaseModel
	Version             string               `gorm:"type:varchar(255)" json:"version"`
	Description         string               `gorm:"type:text" json:"description"`
	ControllerInstances []ControllerInstance `json:"controller_instances"`
}

type Measurement struct {
	CarSessionID uint `json:"car_session_id"`
	SensorID     uint `json:"sensor_id"`
	latency      int
	CreatedAt    *time.Time `gorm:"type:timestamptz" json:"created_at"`
	Data1        float32    `gorm:"type:double precision" json:"data1"`
	Data2        float32    `gorm:"type:double precision" json:"data2"`
	Sensor       Sensor     `json:"sensor"`
	CarSession   CarSession `json:"car_session"`
}

type Sensor struct {
	BaseModel
	ControllerInstanceID uint               `json:"controller_instance_id"`
	Name                 string             `gorm:"type:varchar(255)" json:"name"`
	SensorType           SensorType         `gorm:"type:varchar(255)" json:"sensor_type"`
	ControllerInstance   ControllerInstance `json:"controller_instance"`
}

type Session struct {
	BaseModel
	Name      string     `gorm:"type:varchar(255)" json:"name"`
	StartedAt *time.Time `gorm:"type:timestamptz" json:"started_at"`
	EndedAt   *time.Time `gorm:"type:timestamptz" json:"ended_at"`
	Cars      []Car      `gorm:"many2many:car_sessions;" json:"cars"`
}

type CarSession struct {
	BaseModel
	CarID               uint                 `json:"car_id"`
	SessionID           uint                 `json:"session_id"`
	IsControlledByUser  bool                 `json:"is_controlled_by_user"`
	ControllerInstances []ControllerInstance `gorm:"many2many:car_session_controllers;" json:"controller_instances"`
	Meassurements       []Measurement        `json:"meassurements"`
	Session             Session              `json:"session"`
}

type CarSessionController struct {
	BaseModel
	CarSessionID         uint `json:"car_session_id"`
	ControllerInstanceID uint `json:"controller_instance_id"`
}

type CarController struct {
	BaseModel
	CarID                uint `json:"car_id"`
	ControllerInstanceID uint `json:"controller_instance_id"`
}

type ControllerInstance struct {
	BaseModel
	FirmwareID   uint       `json:"firmware_id"`
	ControllerID uint       `json:"controller_id"`
	name         string     `gorm:"type:varchar(255)" json:"name"`
	Firmware     Firmware   `json:"firmware"`
	Controller   Controller `json:"controller"`
}
