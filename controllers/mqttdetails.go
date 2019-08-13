package controllers

import (
	"github.com/arvindh123/Mod2Mqtt/models"
	"github.com/gin-gonic/gin"
)

func CreateMqtt(c *gin.Context) {
	var mq models.MqttDetails
	db := models.GetDB()
	id := c.Params.ByName("id")
	if id == "0" {
		if err := c.ShouldBindJSON(&mq); err != nil {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
			// c.JSON(200, mq)
		} else {
			if err:= db.Create(&mq).Error; err!= nil {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
			}else {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})
			}
			
		}
	} else if err := db.Where("id = ?", id).First(&mq).Error; err != nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
		
	} else {
		if err = c.ShouldBindJSON(&mq); err != nil {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
			// c.JSON(200, mq)
		} else {
			if err:= db.Save(&mq).Error; err!= nil {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
			}else {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})
			}
		}
	}
}

func GetAllMqtt(c *gin.Context) {
	var mq []models.MqttDetails
	db := models.GetDB()
	// id := c.Params.ByName("id")
	if err := db.Find(&mq).Error; err != nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{err.Error()}}})
	} else {
		c.JSON(200, gin.H{"msg": mq})
	}
}

func DeleteMqtt(c *gin.Context) {
	var mq models.MqttDetails
	db := models.GetDB()
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).Delete(&mq).Error; err == nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{"Done"}}})
	} else {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{err.Error()}}})
	}

}
