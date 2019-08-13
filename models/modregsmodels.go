package models

import "errors"

type ModbusRegisters struct {
	ID             int    `gorm:"UNIQUE;NOT NULL;AUTO_INCREMENT" json:"id"`
	Name           string `gorm:"NOT NULL" json:"name"`
	FunctCode      int    `gorm:"NOT NULL" json:"functcode"`
	Register       uint16 `gorm:"NOT NULL" json:"register"`
	Qty            uint16 `gorm:"NOT NULL" json:"qty"`
	DataType       int    `gorm:"NOT NULL" json:"datatype"`
	ByteOrder      uint8  `gorm:"NOT NULL" json:"byteorder"`
	PostProcess    string `gorm:"NOT NULL" json:"postprocess"`
	Tags           string `gorm:"NOT NULL" json:"tags"`
	DeviceModelsID int    `gorm:"NOT NULL" json:"devicemodelsid" sql:"type:int REFERENCES device_models(id)"`
}

type ModbusRegistersWithModel struct {
	ModbusRegisters
	Model DeviceModels
}

func (mr *ModbusRegisters) BeforeCreate() (err error) {
	var devmod DeviceModels
	if err := db.Where("id = ?", mr.DeviceModelsID).First(&devmod).Error; err != nil {
		if err.Error() == "record not found" {
			return errors.New("Models is not registered in this Edge Gateway")
		}
		return err
	} else {
		if devmod.ID != mr.DeviceModelsID {
			return errors.New("Models is not registered in this Edge Gateway")
		}
		return nil
	}

}
