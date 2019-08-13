package controllers

import (
	"github.com/arvindh123/Mod2Mqtt/models"
	"github.com/gin-gonic/gin"
)

func CreateDeviceDetails(c *gin.Context) {
	var obj models.DeviceDetails
	db := models.GetDB()
	id := c.Params.ByName("id")
	if id == "0" {
		if err := c.ShouldBindJSON(&obj); err != nil {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
		} else {
			if err := db.Create(&obj).Error; err != nil {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
			} else {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})
			}

		}
	} else if err := db.Where("id = ?", id).First(&obj).Error; err != nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
	} else {
		if err = c.ShouldBindJSON(&obj); err != nil {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})

		} else {
			if err := db.Save(&obj).Error; err != nil {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
			} else {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})
			}

		}
	}
}

func GetAllDeviceDetails(c *gin.Context) {
	var dev []models.DeviceDetails
	var obj []models.DeviceAllDetails
	var model models.DeviceModels
	var modregs models.ModbusRegisters
	var inter models.InterfaceDetails
	db := models.GetDB()
	// id := c.Params.ByName("id")
	if err := db.Find(&dev).Error; err != nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{err.Error()}}})
	} else {
		for i := range dev {
			db.Where("id = ?", dev[i].DeviceModelsID).First(&model)
			db.Where("device_models_id = ?", dev[i].DeviceModelsID).Find(&modregs)
			db.Where("id = ?", dev[i].InterfaceDetailsID).First(&inter)
			obj = append(obj, models.DeviceAllDetails{Device: dev[i], Model: model, Interface: inter})
		}
		c.JSON(200, gin.H{"msg": obj})
	}
}

func GetAllDeviceDetailsAllDetails(c *gin.Context) {
	var obj []models.DeviceDetails
	db := models.GetDB()
	// id := c.Params.ByName("id")
	if err := db.Find(&obj).Error; err != nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{err.Error()}}})
	} else {
		c.JSON(200, gin.H{"msg": obj})
	}
}
func DeleteDeviceDetails(c *gin.Context) {
	var obj models.DeviceDetails
	db := models.GetDB()
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).Delete(&obj).Error; err == nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{"Done"}}})
	} else {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{err.Error()}}})
	}

}
