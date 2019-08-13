package controllers

import (
	"strings"

	"github.com/arvindh123/Mod2Mqtt/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	db := models.GetDB()
	id := c.Params.ByName("id")
	if id == "users" {
		GetUsers(c)
	} else {
		var user models.User
		if err := db.Where("id = ?", id).First(&user).Error; err != nil {
			c.JSON(200, []models.Content{models.Content{Content: "User doesn't exist"}})
		} else {
			c.JSON(200, gin.H{"msg": []models.UserDetails{models.UserDetails{ID: user.ID, UserName: user.UserName, Name: user.Name}}})
		}
	}
}

func GetUsers(c *gin.Context) {
	db := models.GetDB()
	var users []models.User
	if err := db.Select("id,user_name, name").Find(&users).Error; err != nil {
		c.JSON(200, gin.H{"msg": []models.Empty{models.Empty{}}})
	} else {
		c.JSON(200, gin.H{"msg": users})
	}
}

func DeleteUser(c *gin.Context) {
	db := models.GetDB()
	id := c.Params.ByName("id")
	var user models.User
	d := db.Where("id = ?", id).Delete(&user)
	_ = d
	c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})
}

func UpdateUser(c *gin.Context) {
	db := models.GetDB()
	var user models.User
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(200, gin.H{"msg": []models.Empty{models.Empty{}}})
	}
	c.BindJSON(&user)
	if err := db.Save(&user).Error; err != nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})
	} else {
		c.JSON(200, gin.H{"msg": []models.User{user}})
	}

}

func GetProfile(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)
	userId := session.Get("user")
	if userId == nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Not logged in"}}})
	} else {
		var user models.User
		if err := db.Where("user_name = ?", userId).First(&user).Error; err != nil {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "User doesn't exist"}}})
		} else {
			c.JSON(200, gin.H{"msg": []models.UserDetails{models.UserDetails{ID: user.ID, Name: user.UserName, UserName: user.Name}}})
		}
	}
}

func UpdateProfile(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)
	userId := session.Get("user")
	if userId == nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Not logged in"}}})
	} else {
		var updateUser models.UserUpdate
		c.Bind(&updateUser)

		var user models.User
		if err := db.Where("user_name = ?", userId).First(&user).Error; err != nil {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "User doesn't exist"}}})
		} else {

			if strings.Trim(updateUser.Password, " ") != "" && updateUser.Password != "" {
				user.Password = updateUser.Password

			}
			if strings.Trim(updateUser.Name, " ") != "" && updateUser.Name != "" {
				user.Name = updateUser.Name

			}
			if err := db.Save(&user).Error; err != nil {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})
			} else {
				c.JSON(200, gin.H{"msg": []models.User{user}})
			}
		}
	}

}

func CreateUser(c *gin.Context) {
	db := models.GetDB()
	var user models.User
	c.BindJSON(&user)

	if user.UserName == "" || user.Password == "" || user.Name == "" {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Enter All Details"}}})

	} else {
		if err := db.Create(&user).Error; err != nil {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
		} else {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})
		}
	}

}

// func PostUpdateProfile(c *gin.Context) {
// 	session := sessions.Default(c)
// 	userId := session.Get("user")
// 	if userId == nil {
// 		c.JSON(200, gin.H{"msg": "Not logged in"})
// 	} else {
// 		var updateUser UserUpdate
// 		c.Bind(&updateUser)

// 		var user User
// 		if err := db.Where("user_name = ?", userId).First(&user).Error; err != nil {
// 			c.JSON(200, gin.H{"msg": err})
// 		} else {
// 			fmt.Println(updateUser.Name, updateUser.Password)
// 			if strings.Trim(updateUser.Password, " ") != "" && updateUser.Password != "" {
// 				user.Password = updateUser.Password
// 				fmt.Println(user.Password)
// 			}
// 			if strings.Trim(updateUser.Name, " ") != "" && updateUser.Name != "" {
// 				user.Name = updateUser.Name
// 				fmt.Println(user.Name)
// 			}
// 			if err:=db.Save(&user).Error; err!=nil {
// 	c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})
// 	} else{
// 		c.JSON(200, gin.H{"msg": []models.User{user}})
// }
// 		}
// 	}

// }
