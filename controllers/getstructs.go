package controllers

import (
	"fmt"

	"github.com/arvindh123/Mod2Mqtt/models"
	"github.com/gin-gonic/gin"
)

// func (d *DeviceWithRegs)AddModRegs(modreg models.ModbusRegisters) []models.ModbusRegisters {
// 	d.ModRegs = append(d.Modregs,modreg)
// 	return d.Modregs
// }

// func (i *InterWithDevices)AddDevRegs(devwithregs DevicesWithRegs) []DevicesWithRegs {
// 	i.Devices = append(i.Devices,devwithregs)
// 	return i.Devices
// }

// func (t *ToRead)AddInterDev(interdev InterWithDevices) []InterWithDevices {
// 	t.ToRead = append(t.ToRead,interdev)
// 	return t.ToRead
// }

func Test(c *gin.Context) {
	var inters []models.InterfaceDetails
	var intersdevice []models.InterWithDevices
	db := models.GetDB()
	var read models.AllStructParams
	if err := db.Find(&inters).Error; err != nil {
		fmt.Println(err.Error())
	} else {
		for i, inter := range inters {
			_ = i
			var devices []models.DeviceDetails
			var devregs []models.DeviceWithRegs

			if err := db.Where("interface_details_id = ?", inter.ID).Find(&devices).Error; err != nil {
				fmt.Println(err.Error())
			} else {

				for j, dev := range devices {
					_ = j
					var modregs []models.ModbusRegisters
					if err := db.Where("device_models_id = ?", dev.DeviceModelsID).Find(&modregs).Error; err != nil {
						fmt.Println(err.Error())
					} else {
						devregs = append(devregs, models.DeviceWithRegs{Device: dev, Modregs: modregs})
					}
				}
			}
			intersdevice = append(intersdevice, models.InterWithDevices{Inter: inter, Devices: devregs})
		}
	}
	read.AllStructParams = intersdevice
	c.JSON(200, gin.H{"msg": read})
}

func GetAllInters() models.AllStructParams {
	var inters []models.InterfaceDetails
	var intersdevice []models.InterWithDevices
	db := models.GetDB()
	var read models.AllStructParams
	if err := db.Find(&inters).Error; err != nil {
		fmt.Println(err.Error())
	} else {
		for i, inter := range inters {
			_ = i
			var devices []models.DeviceDetails
			var devregs []models.DeviceWithRegs

			if err := db.Where("interface_details_id = ?", inter.ID).Find(&devices).Error; err != nil {
				fmt.Println(err.Error())
			} else {

				for j, dev := range devices {
					_ = j
					var modregs []models.ModbusRegisters
					if err := db.Where("device_models_id = ?", dev.DeviceModelsID).Find(&modregs).Error; err != nil {
						fmt.Println(err.Error())
					} else {
						devregs = append(devregs, models.DeviceWithRegs{Device: dev, Modregs: modregs})
					}
				}
			}
			intersdevice = append(intersdevice, models.InterWithDevices{Inter: inter, Devices: devregs})
		}
	}
	read.AllStructParams = intersdevice
	return read
}

func GetMq() (models.MqttDetails, error) {
	var mq models.MqttDetails
	db := models.GetDB()

	if err := db.First(&mq).Error; err != nil {
		return models.MqttDetails{}, err
	}
	return mq,nil
}
