package controllers

import (
	"github.com/arvindh123/Mod2Mqtt/models"
	"github.com/gin-gonic/gin"
)

func CreateModRegs(c *gin.Context) {
	var obj models.ModbusRegisters
	db := models.GetDB()
	id := c.Params.ByName("id")
	if id == "0" {
		if err := c.ShouldBindJSON(&obj); err != nil {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
		} else {
			if err := db.Create(&obj).Error; err != nil {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
				c.Abort()
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

func GetAllModRegs(c *gin.Context) {
	var obj []models.ModbusRegisters
	var finalObj []models.ModbusRegistersWithModel
	db := models.GetDB()
	// id := c.Params.ByName("id")
	_ = finalObj
	if err := db.Find(&obj).Error; err != nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
	} else {
		for i, o := range obj {
			model := models.DeviceModels{}
			_ = i
			db.Where("id = ?", o.DeviceModelsID).First(&model)
		}
		c.JSON(200, gin.H{"msg": obj})
	}
}

func DeleteModRegs(c *gin.Context) {
	var obj models.ModbusRegisters
	db := models.GetDB()
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).Delete(&obj).Error; err == nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})
	} else {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
	}

}
