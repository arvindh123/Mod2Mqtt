package models

import (
	"errors"
)

type DeviceDetails struct {
	ID                 int    `gorm:"UNIQUE;NOT NULL;AUTO_INCREMENT" json:"id"`
	Name               string `gorm:"NOT NULL" json:"name"`
	DeviceId           string `gorm:"UNIQUE;NOT NULL" json:"deviceid"`
	MbId               byte   `gorm:"NOT NULL" json:"mbid"`
	DeviceModelsID     int    `gorm:"NOT NULL" json:"devicemodelsid" sql:"type:int REFERENCES device_models(id) ON DELETE RESTRICT"`
	InterfaceDetailsID int    `gorm:"NOT NULL" json:"intefacedetailsid" sql:"type:int REFERENCES interface_details(id) ON DELETE RESTRICT"`
}

type DeviceAllDetails struct {
	Device    DeviceDetails
	Model     DeviceModels
	Interface InterfaceDetails
}
type DeviceModels struct {
	ID    int    `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	Make  string `gorm:"NOT NULL" json:"make"`
	Model string `gorm:"UNIQUE;NOT NULL" json:"model"`
}

type DeviceModelsRegs struct {
	Model   DeviceModels
	ModRegs []ModbusRegisters
}

func (dev *DeviceDetails) BeforeCreate() (err error) {
	var devmod DeviceModels
	if err := db.Where("id = ?", dev.DeviceModelsID).First(&devmod).Error; err != nil {
		if err.Error() != "record not found" {
			return errors.New("Model not Supported in this Device, To add support Contact Admin Team")
		}
		return err
	} else {
		if devmod.ID != dev.DeviceModelsID {
			return errors.New("Model not Supported in this Device, To add support Contact Admin Team")
		}
		return nil
	}

}

func (modls *DeviceModels) BeforeDelete() (err error) {
	var devices []DeviceDetails
	db.Where("device_models_id = ?", modls.ID).Find(&devices)
	if err := db.Where("device_models_id = ?", modls.ID).Find(&devices).Error; err != nil {
		if err.Error() != "record not found" {
			//----------------------------------------------------------------------------------//
			var modregs []ModbusRegisters
			if err := db.Where("device_models_id = ?", modls.ID).Find(&modregs).Error; err != nil {
				if err.Error() != "record not found" {
					return nil
				}
				return err
			}
			if len(modregs) <= 0 {
				return nil
			}
			return errors.New("Model Can't delete until related Modbus Registers are removed")
			//----------------------------------------------------------------------------------//
		}
		return err
	}
	if len(devices) <= 0 {
		//----------------------------------------------------------------------------------//
		var modregs []ModbusRegisters
		if err := db.Where("device_models_id = ?", modls.ID).Find(&modregs).Error; err != nil {
			if err.Error() != "record not found" {
				return nil
			}
			return err
		}
		if len(modregs) <= 0 {
			return nil
		}
		return errors.New("Model Can't delete until related Modbus Registers are removed")
		//----------------------------------------------------------------------------------//
	}
	return errors.New("Model Can't delete until related Device removed")
}
