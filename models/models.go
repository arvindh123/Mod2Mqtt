package models

// import (
// 	"github.com/jinzhu/gorm"
// )

type Mqtt struct {
	Id       int    `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY" json:"id"`
	Ip       string `gorm:"UNIQUE;NOT NULL" json:"ip"`
	Port     int    `gorm:"NOT NULL" json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Topics struct {
	Id      int                `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY" json:"id"`
	Topic   string             `gorm:"NOT NULL" json:"topic"`
	Qos     int                `gorm:"NOT NULL" json:"qos"`
	Retain  bool               `gorm:"NOT NULL" json:"retain"`
	Delay   int                `gorm:"NOT NULL" json:"delay"`
	ModRegs []*ModbusRegisters `gorm:"many2many:regs_topics;"  json:"modregs"`
}
type SerialDetails struct {
	Id       int                `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY" json:"id"`
	ComPort  string             `gorm:"UNIQUE;NOT NULL" json:"comport"`
	BaudRate int                `gorm:"NOT NULL" json:"baudrate"`
	DataBits int                `gorm:"NOT NULL" json:"databits"`
	Parity   string             `gorm:"NOT NULL;size:1" json:"parity"`
	StopBits int                `gorm:"NOT NULL" json:"stopbits"`
	Timeout  int                `gorm:"NOT NULL" json:"timeout"`
	ModRegs  []*ModbusRegisters `gorm:"many2many:regs_ports;"  json:"modregs"`
}

type ModbusRegisters struct {
	ID          int              `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	Name        string           `gorm:"NOT NULL" json:"name"`
	Unit        byte             `gorm:"NOT NULL" json:"unit"`
	FunctCode   int              `gorm:"NOT NULL" json:"functcode"`
	Register    uint16           `gorm:"NOT NULL" json:"register"`
	Qty         uint16           `gorm:"NOT NULL" json:"qty"`
	DataType    int              `gorm:"NOT NULL" json:"datatype"`
	ByteOrder   uint8            `gorm:"NOT NULL" json:"byteorder"`
	PostProcess string           `gorm:"NOT NULL" json:"postprocess"`
	Tags        string           `gorm:"NOT NULL" json:"tags"`
	MqTopic     []*Topics        `gorm:"many2many:regs_topics;" json:"mqtopic"`
	SerialPorts []*SerialDetails `gorm:"many2many:regs_ports;" json:"serialport"`
}
