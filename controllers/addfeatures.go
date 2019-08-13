package controllers

import (
	"github.com/arvindh123/Mod2Mqtt/models"
	"github.com/gin-gonic/gin"
)

func CreateAddFeatures(c *gin.Context) {
	db := models.GetDB()
	var obj models.AddonFeatures
	id := c.Params.ByName("id")
	if id == "0" {
		if err := c.ShouldBindJSON(&obj); err != nil {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
			c.Abort()
		} else {

			if err := db.Create(&obj).Error; err != nil {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
				c.Abort()
			} else {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{"Done"}}})
			}
		}
	} else {
		if err := db.Where("id = ?", id).First(&obj).Error; err != nil {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
		} else {
			if err := c.ShouldBindJSON(&obj); err != nil {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
				c.Abort()
			} else {
				if err := db.Save(&obj).Error; err != nil {
					c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
				} else {
					c.JSON(200, gin.H{"msg": []models.Content{models.Content{"Done"}}})
				}

			}
		}
	}

}

func GetAllAddFeatures(c *gin.Context) {
	db := models.GetDB()
	var obj []models.AddonFeatures
	if err := db.Find(&obj).Error; err != nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
	} else {
		// for i, port := range obj {
		// 	var modregs []*ModbusRegisters
		// 	if err = db.Model(&port).Association("ModRegs").Find(&modregs).Error; err != nil {
		// 		// topics[i].ModRegs = modregs
		// 	} else {
		// 		obj[i].ModRegs = modregs
		// 	}
		// }
		c.JSON(200, gin.H{"msg": obj})
	}
}

func DeleteAddFeatures(c *gin.Context) {
	db := models.GetDB()
	id := c.Params.ByName("id")
	var obj models.AddonFeatures
	if err := db.Where("id = ?", id).Delete(&obj).Error; err != nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
	} else {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})
	}

}
