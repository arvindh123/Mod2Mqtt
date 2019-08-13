package controllers

import (
	"github.com/arvindh123/Mod2Mqtt/models"
	"github.com/gin-gonic/gin"
)

func CreateDeviceModels(c *gin.Context) {
	var obj models.DeviceModels
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

func GetAllDeviceModels(c *gin.Context) {
	var obj []models.DeviceModels
	db := models.GetDB()
	// id := c.Params.ByName("id")
	if err := db.Find(&obj).Error; err != nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{err.Error()}}})
	} else {
		c.JSON(200, gin.H{"msg": obj})
	}
}

func DeleteDeviceModels(c *gin.Context) {
	var obj models.DeviceModels
	db := models.GetDB()
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&obj).Error; err == nil {
		if err := db.Delete(&obj).Error; err == nil {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{"Done"}}})
		} else {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{err.Error()}}})
		}

	} else {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{err.Error()}}})
	}

}
