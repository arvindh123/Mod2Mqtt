package main

// only need mysql OR sqlite
// both are included here for reference

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/gin-contrib/pprof"

	// "github.com/gin-contrib/pprof"
	"github.com/arvindh123/Mod2Mqtt/controllers"
	"github.com/arvindh123/Mod2Mqtt/models"
	"github.com/arvindh123/Mod2Mqtt/proc"
	"github.com/arvindh123/Mod2Mqtt/ws"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

var wsClients = make(map[*websocket.Conn]bool) // connected clients
var wsClientsChan = make(chan map[*websocket.Conn]bool)

var broadcast = make(chan rdmsg) // broadcast channel

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type rdmsg struct {
	ty  int
	msg []byte
}

func main() {

	models.SetDB("sqlite3", "./gorm.db")
	models.AutoMigrate()
	// models.SetDB("mysql", "root:root@tcp(127.0.0.1:3306)/schema_name")
	if err != nil {
		fmt.Println(err)
	}
	db := models.GetDB()
	defer db.Close()

	r := gin.Default()
	pprof.Register(r)
	r.Use(static.Serve("/local/", static.LocalFile("./views", true)))
	Store := sessions.NewCookieStore([]byte("anosecret"))

	r.Use(sessions.Sessions("mysession", Store))

	r.Use(func(c *gin.Context) {
		defer func() {
			session := sessions.Default(c)
			session.Clear()
			session.Save()
		}()
		c.Next()
	})
	r.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "welcome"})
	})

	r.GET("/cookie", controllers.SetCookies)

	r.POST("/login", controllers.Login)
	r.GET("/logout", controllers.Logout)

	r.POST("/FirstTimeCreate", controllers.CreateUser)

	FirstTime := r.Group("/FirstTime")
	FirstTime.Use(controllers.CheckforNoUsers())
	{
		FirstTime.POST("/create", controllers.CreateUser)
	}

	user := r.Group("/user")
	user.Use(controllers.AuthRequired())
	{
		user.POST("/create", controllers.CreateUser)
		user.GET("/user/:id", controllers.GetUser)
		user.GET("/profile", controllers.GetProfile)
		user.GET("/users", controllers.GetUsers)
		user.PUT("/update", controllers.UpdateProfile)
		// user.PUT("/update/:id", controllers.UpdateUser) // need to add super user validation then only we can implement this feature
		user.DELETE("/delete/:id", controllers.DeleteUser)
	}

	apiV1 := r.Group("/api/v1/")
	apiV1.Use(controllers.AuthRequired())
	{
		apiV1.POST("/interface/create/:id", controllers.CreateMbInterface)
		apiV1.GET("/interface/getall", controllers.GetAllMbInterface)
		apiV1.DELETE("/interface/delete/:id", controllers.DeleteMbInterface)

		apiV1.POST("/mqtt/create/:id", controllers.CreateMqtt)
		apiV1.GET("/mqtt/getall", controllers.GetAllMqtt)
		apiV1.DELETE("/mqtt/delete/:id", controllers.DeleteMqtt)

		apiV1.POST("/devicemodels/create/:id", controllers.CreateDeviceModels)
		apiV1.GET("/devicemodels/getall", controllers.GetAllDeviceModels)
		apiV1.DELETE("/devicemodels/delete/:id", controllers.DeleteDeviceModels)

		apiV1.POST("/modregs/create/:id", controllers.CreateModRegs)
		apiV1.GET("/modregs/getall", controllers.GetAllModRegs)
		apiV1.DELETE("/modregs/delete/:id", controllers.DeleteModRegs)

		apiV1.POST("/devices/create/:id", controllers.CreateDeviceDetails)
		apiV1.GET("/devices/getall", controllers.GetAllDeviceDetails)
		apiV1.DELETE("/devices/delete/:id", controllers.DeleteDeviceDetails)

		apiV1.POST("/addfeatures/create/:id", controllers.CreateAddFeatures)
		apiV1.GET("/addfeatures/getall", controllers.GetAllAddFeatures)
		apiV1.DELETE("/addfeatures/delete/:id", controllers.DeleteAddFeatures)

		apiV1.GET("/test", controllers.Test)

	}
	proc.ManInit()
	start := make(chan bool)
	stop := make(chan bool)
	go handleMessages(start, stop)
	go proc.TryToStart(start, stop)
	r.Run(":5000")
	// log.Println("serving on port 8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
	os.Exit(1)

	// mod.ModStart()

}

func wshandler(w http.ResponseWriter, r *http.Request) {

	wsupgrader.CheckOrigin = func(r *http.Request) bool {
		fmt.Println("res - %v", r)
		return true
	}
	conn, err := wsupgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Failed to upgrade ws: %+v", err)
		return
	}
	wsClients[conn] = true
	ws.WsClients = wsClients

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		broadcast <- rdmsg{t, msg}
		// conn.WriteMessage(t, msg)
	}

}

func handleMessages(start, stop chan bool) {
	var dat map[string]interface{}
	for {
		msg := <-broadcast
		// fmt.Println("Msg  - ", msg)
		// fmt.Println("Msg  - ", msg.msg)
		if err := json.Unmarshal(msg.msg, &dat); err != nil {
			fmt.Println("Error in json Unmarshal", err)
		} else {
			for k, v := range dat {
				if k == "Cmd" {
					if v == "1" {
						select {
						case start <- true:
						}
					} else if v == "2" {
						select {
						case stop <- true:
						}
					} else if v == "3" {
						select {
						case stop <- true:
						}
					}
				} else {
					for client := range wsClients {
						err := client.WriteMessage(msg.ty, msg.msg)
						if err != nil {
							// log.Printf("error: %v", err)
							client.Close()
							delete(wsClients, client)
						}
					}
				}
			}
		}
	}
}
