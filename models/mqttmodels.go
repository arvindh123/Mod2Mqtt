package models

type MqttDetails struct {
	ID       int    `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	Ip       string `gorm:"UNIQUE;NOT NULL" json:"ip"`
	Port     int    `gorm:"NOT NULL" json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
