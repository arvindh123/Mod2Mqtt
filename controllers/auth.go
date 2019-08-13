package controllers

import (
	"net/http"
	"strings"

	"github.com/arvindh123/Mod2Mqtt/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SetCookies(c *gin.Context) {
	db := models.GetDB()
	session := sessions.Default(c)
	userId := session.Get("user")

	if userId == nil {
		c.SetCookie("Gate_user", "{\"loggedin\":\"false\",\"username\":\"anonymous\",\"name\":\"anonymous\",\"level\":\"0\"}", 3600, "/", "localhost", false, false)

	} else {
		var user models.User
		if err := db.Where("user_name = ?", userId).First(&user).Error; err != nil {
			c.SetCookie("Gate_user", "{\"loggedin\":\"false\",\"username\":\"anonymous\",\"name\":\"anonymous\",\"level\":\"0\"}", 3600, "/", "localhost", false, false)
		} else {
			c.SetCookie("Gate_user", "{\"loggedin\":\"true\",\"username\":\""+user.UserName+"\",\"name\":\""+user.UserName+"\",\"level\":\"1\"}", 3600, "/", "localhost", false, false)
		}

	}

}

func CheckforNoUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := models.GetDB()
		var user models.User

		if err := db.First(&user).Error; err != nil {
			if user != (models.User{}) {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: err.Error()}}})
				c.Abort()
			} else {
				c.Next()
			}
		} else {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "First User already created"}}})
			c.Abort()
		}
	}
}
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			// You\"d normally redirect to login page

			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Not logged in"}}})
			SetCookies(c)
			c.Abort()
		} else {
			// Continue down the chain to handler etc
			// c.JSON(200, gin.H{"loggedIn": "true"})
			c.Next()
		}
	}
}

func Login(c *gin.Context) {
	db := models.GetDB()
	var login models.UserLogin
	var user models.User
	session := sessions.Default(c)
	c.BindJSON(&login)
	if strings.Trim(login.UserName, " ") == "" || strings.Trim(login.Password, " ") == "" {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Enter All details"}}})
		return
	}
	if err := db.Where("user_name = ?", login.UserName).First(&user).Error; err != nil {
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Username \"" + login.UserName + "\" does not exists"}}})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err == nil {
		session.Set("user", user.UserName)
		err := session.Save()
		if err != nil {
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Failed to generate session token"}}})
			return
		} else {

			SetCookies(c)
			c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})
			return
		}
	} else {
		if user.Password == login.Password {
			session.Set("user", user.UserName)
			err := session.Save()
			if err != nil {
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Failed to generate session token"}}})
				return
			} else {
				SetCookies(c)
				c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})
				return
			}
		}
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Password Wrong"}}})
		return
	}
}

func Logout(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		SetCookies(c)
		c.JSON(200, gin.H{"msg": []models.Content{models.Content{Content: "Not logged in"}}})
		SetCookies(c)
	} else {
		session.Delete("user")
		session.Save()
		SetCookies(c)
		c.JSON(http.StatusOK, gin.H{"msg": []models.Content{models.Content{Content: "Done"}}})

	}
}
