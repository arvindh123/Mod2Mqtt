package models

import "github.com/jinzhu/gorm"

var db *gorm.DB

//SetDB establishes connection to database and saves its handler into db *sqlx.DB
func SetDB(typ, connection string) {
	var err error
	db, err = gorm.Open(typ, connection)
	if err != nil {
		panic(err)
	}
}

//GetDB returns database handler
func GetDB() *gorm.DB {
	return db
}

//AutoMigrate runs gorm auto migration
func AutoMigrate() {
	db.Debug().AutoMigrate(&User{}, &InterfaceDetails{}, &MqttDetails{}, &DeviceDetails{}, &DeviceModels{}, &ModbusRegisters{}, &AddonFeatures{})
}

type Empty struct {
}

type Content struct {
	Content string `json:"content"`
}
