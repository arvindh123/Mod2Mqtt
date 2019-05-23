package main

// only need mysql OR sqlite
// both are included here for reference

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"

	"./mod"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/goburrow/modbus"
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

type MyRtuHandler struct {
	*modbus.RTUClientHandler
}

type MyTcpHandler struct {
	*modbus.TCPClientHandler
}

type rdmsg struct {
	ty  int
	msg []byte
}

type UserLogin struct {
	UserName string `gorm:"UNIQUE;NOT NULL" json:"username"`
	Password string `gorm:"NOT NULL" json:"password"`
}
type Empty struct {
}

type AssoPortsStruct struct {
	PortId []*int `json:"portids"`
}

type AssoRegsStruct struct {
	ModRegId []*int `json:"modregids"`
}

type AssoTopicsStruct struct {
	TopicsId []*int `json:"topicids"`
}
type User struct {
	ID   int    `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	Name string `gorm:"NOT NULL" json:"name"`
	UserLogin
}

type UserDetails struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	UserName string `json:"username"`
}
type Content struct {
	Content string `json:"content"`
}

type UserUpdate struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Mqtt struct {
	Id       int    `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY" json:"id"`
	Ip       string `gorm:"UNIQUE;NOT NULL" json:"ip"`
	Port     int    `gorm:"NOT NULL" json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Topics struct {
	Id      int                `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY" json:"id"`
	Topic   string             `gorm:"NOT NULL" json:"topic"`
	Qos     int                `gorm:"NOT NULL" json:"qos"`
	Retain  bool               `gorm:"NOT NULL" json:"retain"`
	Delay   int                `gorm:"NOT NULL" json:"delay"`
	ModRegs []*ModbusRegisters `gorm:"many2many:regs_topics;"  json:"modregs"`
}

type SerialDetails struct {
	Id       int                `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY" json:"id"`
	Name     string             `gorm:"NOT NULL" json:"name"`
	Type     int                `gorm:"default:0" json:"type"`
	IpAdd    string             `gorm:"default:''" json:"ipadd"`
	Port     int                `gorm:"default:0" json:"port"`
	ComPort  string             `gorm:"default:''" json:"comport"`
	BaudRate int                `gorm:"default:0" json:"baudrate"`
	DataBits int                `gorm:"default:0" json:"databits"`
	Parity   string             `gorm:"default:''" gorm:"size:1" json:"parity"`
	StopBits int                `gorm:"default:0" json:"stopbits"`
	Timeout  int                `gorm:"default:0" json:"timeout"`
	ModRegs  []*ModbusRegisters `gorm:"many2many:regs_ports;"  json:"modregs"`
}

type ModbusRegisters struct {
	ID          int              `gorm:"UNIQUE;NOT NULL;PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	Name        string           `gorm:"NOT NULL" json:"name"`
	Unit        byte             `gorm:"NOT NULL" json:"unit"`
	FunctCode   int              `gorm:"NOT NULL" json:"functcode"`
	Register    uint16           `gorm:"NOT NULL" json:"register"`
	Qty         uint16           `gorm:"NOT NULL" json:"qty"`
	DataType    int              `gorm:"NOT NULL" json:"datatype"`
	ByteOrder   uint8            `gorm:"NOT NULL" json:"byteorder"`
	PostProcess string           `gorm:"NOT NULL" json:"postprocess"`
	Tags        string           `gorm:"NOT NULL" json:"tags"`
	MqTopic     []*Topics        `gorm:"many2many:regs_topics;" json:"mqtopic"`
	SerialPorts []*SerialDetails `gorm:"many2many:regs_ports;" json:"serialport"`
}

func main() {
	// NOTE: See we’re using = to assign the global var
	// instead of := which would assign it only in this function
	var RtuHandlers []*modbus.RTUClientHandler
	var RtuErrs []error
	var TcpHandlers []*modbus.TCPClientHandler
	var TcpErrs []error
	var mqErr error
	var mqClient mqtt.Client
	var wg sync.WaitGroup
	SpanStopper := make(map[int](chan int))
	PortChance := make(map[string](chan int))
	var err error
	var status bool = false

	webModMqChan := make(chan int)
	_ = webModMqChan
	db, err = gorm.Open("sqlite3", "./gorm.db")
	// db, _ = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_schema")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	// db.AutoMigrate(&User{})
	db.Debug().AutoMigrate(&User{}, &SerialDetails{}, &ModbusRegisters{}, &Mqtt{}, &Topics{})

	// go mod.ModMqProcess(db, webModMqChan, wsClientsChan)

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./views", true)))
	Store := sessions.NewCookieStore([]byte("anosecret"))

	r.Use(sessions.Sessions("mysession", Store))

	r.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "welcome"})
	})

	r.GET("/cookie", setCookies)

	r.POST("/login", Login)
	r.GET("/logout", Logout)

	user := r.Group("/user")
	user.Use(AuthRequired())
	{
		user.POST("/create", CreateUser)
		user.GET("/user/:id", GetUser)
		user.GET("/profile", GetProfile)
		user.GET("/users", GetUsers)
		user.PUT("/update", UpdateProfile)
		// user.PUT("/update/:id", UpdateUser)
		user.DELETE("/delete/:id", DeleteUser)
	}

	apiV1 := r.Group("/api/v1/")
	{
		apiV1.POST("/modbus/params/:id", CreateModbusSerailParams)
		apiV1.GET("/modbus/params", GetModbusSerailParams)
		apiV1.DELETE("/modbus/params/:id", DeleteModbusSerailParams)

		apiV1.POST("/mqtt/params", CreateMqttParams)
		// apiV1.PUT("/mqtt/params", CreateMqttParams)
		apiV1.GET("/mqtt/params", GetMqttParams)

		apiV1.POST("/modbus/regs/:id", CreateModRegs)
		apiV1.GET("/modbus/regs", GetModRegs)
		apiV1.DELETE("/modbus/regs/:id", DeleteModRegs)

		apiV1.POST("/mqtt/topic/:id", CreateTopics)
		apiV1.GET("/mqtt/topics", GetTopics)
		apiV1.DELETE("/mqtt/topics/:id", DeleteTopics)

		apiV1.POST("/topics/modregs/:id", AssoTopic2Regs)
		apiV1.DELETE("/topics/modregs/:id", DelTopic2Regs)
		apiV1.DELETE("/topics/modregs/:id/all", DelTopic2RegsAll)

		apiV1.POST("/modregs/topics/:id", AssoRegs2Topic)
		apiV1.DELETE("/modregs/topics/:id", DelRegs2Topic)

		apiV1.POST("/serial/modregs/:id", AssoPort2Regs)
		apiV1.DELETE("/serial/modregs/:id", DelPort2Regs)
		apiV1.DELETE("/serial/modregs/:id/all", DelPort2RegsAll)
	}
	go handleMessages(status, RtuHandlers, TcpHandlers, mqClient, SpanStopper, PortChance, &wg, RtuErrs, TcpErrs, mqErr)
	r.Run(":5000")
	// log.Println(http.ListenAndServe("localhost:6060", nil))

}

func wshandler(w http.ResponseWriter, r *http.Request) {

	conn, err := wsupgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Failed to upgrade ws: %+v", err)
		return
	}
	wsClients[conn] = true
	mod.WsClients = wsClients

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		broadcast <- rdmsg{t, msg}
		// conn.WriteMessage(t, msg)
	}

}

func handleMessages(status bool, RtuHandlers []*modbus.RTUClientHandler, TcpHandlers []*modbus.TCPClientHandler, mqClient mqtt.Client, SpanStopper map[int](chan int), PortChance map[string](chan int), wg *sync.WaitGroup, RtuErrs []error, TcpErrs []error, mqErr error) {
	var dat map[string]interface{}
	for {
		msg := <-broadcast
		// fmt.Println("Msg type - ", msg.ty)
		// fmt.Println("Msg  - ", msg.msg)
		if err := json.Unmarshal(msg.msg, &dat); err != nil {
			fmt.Println("Error in json Unmarshal", err)
		} else {
			for k, v := range dat {
				if k == "Cmd" {
					if v == "1" {
						if !status {
							status, RtuHandlers, TcpHandlers, mqClient, SpanStopper, PortChance, wg, RtuErrs, TcpErrs, mqErr = mod.MultiModMqProcessStart(db, status, wsClientsChan)
						} else {
							fmt.Println("Already Running")
						}
						if mqErr != nil {
							fmt.Println("Stoping due to MqErr")
							mod.MultiModMqProcessStop(RtuHandlers, TcpHandlers, mqClient, SpanStopper, PortChance, wg)
							status = false
						}

						for _, modErr := range RtuErrs {
							if modErr != nil {
								fmt.Println("Stoping due to Modbus RTU Error")
								mod.MultiModMqProcessStop(RtuHandlers, TcpHandlers, mqClient, SpanStopper, PortChance, wg)
								status = false
							}
						}

						for _, modErr := range TcpErrs {
							if modErr != nil {
								fmt.Println("Stoping due to Modbus TCP Error")
								mod.MultiModMqProcessStop(RtuHandlers, TcpHandlers, mqClient, SpanStopper, PortChance, wg)
								status = false
							}
						}
					} else if v == "2" {
						err = mod.MultiModMqProcessStop(RtuHandlers, TcpHandlers, mqClient, SpanStopper, PortChance, wg)
						status = false
					}
				} else {
					for client := range wsClients {
						err := client.WriteMessage(msg.ty, msg.msg)
						if err != nil {
							log.Printf("error: %v", err)
							client.Close()
							delete(wsClients, client)
						}
					}
				}
			}
		}

	}
}

func DeleteUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user User
	d := db.Where("id = ?", id).Delete(&user)
	_ = d
	c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
}

func UpdateUser(c *gin.Context) {
	var user User
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(200, gin.H{"msg": []Empty{Empty{}}})
	}
	c.BindJSON(&user)
	db.Save(&user)
	c.JSON(200, gin.H{"msg": []User{user}})
}

func GetProfile(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("user")
	if userId == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": []Content{Content{"Not logged in"}}})
	} else {
		var user User
		if err := db.Where("user_name = ?", userId).First(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": []Content{Content{"User doesn't exist"}}})
		} else {
			c.JSON(200, gin.H{"msg": []UserDetails{UserDetails{Id: user.ID, Name: user.UserName, UserName: user.Name}}})
		}
	}
}

func UpdateProfile(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("user")
	if userId == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": []Content{Content{"Not logged in"}}})
	} else {
		var updateUser UserUpdate
		c.Bind(&updateUser)

		var user User
		if err := db.Where("user_name = ?", userId).First(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": []Content{Content{"User doesn't exist"}}})
		} else {

			if strings.Trim(updateUser.Password, " ") != "" && updateUser.Password != "" {
				user.Password = updateUser.Password

			}
			if strings.Trim(updateUser.Name, " ") != "" && updateUser.Name != "" {
				user.Name = updateUser.Name

			}
			db.Save(&user)
			c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
		}
	}

}

func CreateUser(c *gin.Context) {
	var user User
	c.BindJSON(&user)

	if user.UserName == "" || user.Password == "" || user.Name == "" {
		c.JSON(200, gin.H{"msg": []Content{Content{"Enter All Details"}}})
		c.Abort()
	} else {
		if err := db.Create(&user).Error; err != nil {
			c.JSON(200, gin.H{"msg": []Content{Content{"Error in creating user"}}})
		} else {
			c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
		}
	}

}

func setCookies(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get("user")

	if userId == nil {
		c.SetCookie("Gate_user", "{\"loggedin\":\"false\",\"username\":\"anonymous\",\"name\":\"anonymous\",\"level\":\"0\"}", 3600, "/", "localhost", false, false)

	} else {
		var user User
		if err := db.Where("user_name = ?", userId).First(&user).Error; err != nil {
			c.SetCookie("Gate_user", "{\"loggedin\":\"false\",\"username\":\"anonymous\",\"name\":\"anonymous\",\"level\":\"0\"}", 3600, "/", "localhost", false, false)
		} else {
			c.SetCookie("Gate_user", "{\"loggedin\":\"true\",\"username\":\""+user.UserName+"\",\"name\":\""+user.UserName+"\",\"level\":\"1\"}", 3600, "/", "localhost", false, false)
		}

	}

}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	if id == "users" {
		GetUsers(c)
	} else {
		var user User
		if err := db.Where("id = ?", id).First(&user).Error; err != nil {
			c.JSON(200, []Content{Content{"User doesn't exist"}})
		} else {
			c.JSON(200, gin.H{"msg": []UserDetails{UserDetails{Id: user.ID, UserName: user.UserName, Name: user.Name}}})
		}
	}
}

func GetUsers(c *gin.Context) {
	var users []User
	if err := db.Select("id,user_name, name").Find(&users).Error; err != nil {
		c.JSON(200, gin.H{"msg": []Empty{Empty{}}})
	} else {
		c.JSON(200, gin.H{"msg": users})
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			// You\"d normally redirect to login page
			c.JSON(200, gin.H{"msg": []Content{Content{"Not logged in"}}})
			setCookies(c)
			c.Abort()
		} else {
			// Continue down the chain to handler etc
			// c.JSON(200, gin.H{"loggedIn": "true"})
			c.Next()
		}
	}
}

func Login(c *gin.Context) {
	var login UserLogin
	var user User
	session := sessions.Default(c)
	c.BindJSON(&login)
	if strings.Trim(login.UserName, " ") == "" || strings.Trim(login.Password, " ") == "" {
		c.JSON(200, gin.H{"msg": []Content{Content{"Enter All details"}}})
		return
	}
	if err = db.Where("user_name = ?", login.UserName).First(&user).Error; err != nil {
		c.JSON(200, gin.H{"msg": []Content{Content{"Username \"" + login.UserName + "\" does not exists"}}})
		return
	}
	if login.Password == user.Password {
		session.Set("user", user.UserName)
		err := session.Save()
		if err != nil {
			c.JSON(200, gin.H{"msg": []Content{Content{"Failed to generate session token"}}})
			return
		} else {
			setCookies(c)
			c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
			return
		}
	} else {
		c.JSON(200, gin.H{"msg": []Content{Content{"Password Wrong"}}})
		return
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		setCookies(c)
		c.JSON(200, gin.H{"msg": []Content{Content{"Not logged in"}}})
		setCookies(c)
	} else {
		session.Delete("user")
		session.Save()
		setCookies(c)
		c.JSON(http.StatusOK, gin.H{"msg": []Content{Content{"Done"}}})

	}
}

func CreateModbusSerailParams(c *gin.Context) {

	id := c.Params.ByName("id")
	// var modregs ModbusRegisters
	var moddb SerialDetails
	if id == "0" {
		if err := c.ShouldBindJSON(&moddb); err != nil {
			c.JSON(200, gin.H{"msg": []Content{Content{"Invalid format or datatype"}}})
			c.Abort()
		} else {

			if err := db.Create(&moddb).Error; err != nil {
				c.JSON(200, gin.H{"msg": []Content{Content{"Invalid format or datatype or Data already exist"}}})
			} else {
				c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
			}
		}
	} else {
		if err := db.Where("id = ?", id).First(&moddb).Error; err != nil {
			c.JSON(200, gin.H{"msg": []Content{Content{"Id doesn't exist"}}})
		} else {
			if err := c.ShouldBindJSON(&moddb); err != nil {
				c.JSON(200, gin.H{"msg": []Content{Content{"Invalid format or datatype"}}})
				c.Abort()
			} else {
				db.Save(&moddb)
				c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
			}
		}
	}
}

func GetModbusSerailParams(c *gin.Context) {
	var ports []SerialDetails
	if err := db.Find(&ports).Error; err != nil {
		c.JSON(200, gin.H{"msg": []Empty{Empty{}}})
	} else {
		for i, port := range ports {
			var modregs []*ModbusRegisters
			if err = db.Model(&port).Association("ModRegs").Find(&modregs).Error; err != nil {
				// topics[i].ModRegs = modregs
			} else {
				ports[i].ModRegs = modregs
			}
		}
		c.JSON(200, gin.H{"msg": ports})
	}
}

func DeleteModbusSerailParams(c *gin.Context) {
	id := c.Params.ByName("id")
	var mod SerialDetails
	d := db.Where("id = ?", id).Delete(&mod)
	_ = d
	c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
}

func CreateModRegs(c *gin.Context) {
	// var modregs ModbusRegisters
	id := c.Params.ByName("id")
	var modregs ModbusRegisters
	if id == "0" {
		if err := c.ShouldBindJSON(&modregs); err != nil {
			c.JSON(200, gin.H{"msg": []Content{Content{"Invalid format or datatype"}}})
			c.Abort()
		} else {

			if err := db.Create(&modregs).Error; err != nil {
				c.JSON(200, gin.H{"msg": []Content{Content{"Invalid format or datatype"}}})
			} else {
				c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
			}
		}
	} else {
		if err := db.Where("id = ?", id).First(&modregs).Error; err != nil {
			c.JSON(200, gin.H{"msg": []Content{Content{"Id doesn't exist"}}})
		} else {
			if err := c.ShouldBindJSON(&modregs); err != nil {
				c.JSON(200, gin.H{"msg": []Content{Content{"Invalid format or datatype"}}})
				c.Abort()
			} else {
				db.Save(&modregs)
				c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
			}
		}
	}
}

func GetModRegs(c *gin.Context) {
	var modregs []ModbusRegisters
	if err := db.Find(&modregs).Error; err != nil {
		c.JSON(200, gin.H{"msg": []Empty{Empty{}}})
	} else {
		for i, modreg := range modregs {
			var topics []*Topics
			var ports []*SerialDetails
			if err = db.Model(&modreg).Association("MqTopic").Find(&topics).Error; err != nil {

			} else {
				modregs[i].MqTopic = topics
			}
			if err = db.Model(&modreg).Association("SerialPorts").Find(&ports).Error; err != nil {

			} else {
				modregs[i].SerialPorts = ports
			}
		}
		c.JSON(200, gin.H{"msg": modregs})
	}
}

func DeleteModRegs(c *gin.Context) {
	id := c.Params.ByName("id")
	var modreg ModbusRegisters
	d := db.Where("id = ?", id).Delete(&modreg)
	_ = d
	c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})

}

func CreateMqttParams(c *gin.Context) {
	var mq Mqtt
	if err := db.Where("id = ?", 1).First(&mq).Error; err != nil {
		mq.Id = 1
		if err = c.ShouldBindJSON(&mq); err != nil {

			c.JSON(200, gin.H{"msg": []Content{Content{"Invalid format or datatype"}}})
			// c.JSON(200, mq)
			c.Abort()
		} else {
			db.Save(&mq)
			c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
		}
	} else {
		if err = c.ShouldBindJSON(&mq); err != nil {
			c.JSON(200, gin.H{"msg": []Content{Content{"Invalid format or datatype"}}})
			// c.JSON(200, mq)
			c.Abort()
		} else {
			db.Save(&mq)
			c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
		}

	}

}

func GetMqttParams(c *gin.Context) {
	var mq Mqtt
	// id := c.Params.ByName("id")
	if err := db.Where("id = ?", 1).First(&mq).Error; err != nil {
		c.JSON(200, gin.H{"msg": []Empty{Empty{}}})
	} else {
		c.JSON(200, gin.H{"msg": []Mqtt{mq}})
	}
}

func CreateTopics(c *gin.Context) {
	id := c.Params.ByName("id")
	var topic Topics
	if id == "0" {
		if err := c.ShouldBindJSON(&topic); err != nil {
			c.JSON(200, gin.H{"msg": []Content{Content{"Invalid format or datatype"}}})
			c.Abort()
		} else {
			if err := db.Create(&topic).Error; err != nil {
				c.JSON(200, gin.H{"msg": []Content{Content{"Invalid format or datatype"}}})
			} else {
				c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
			}
		}
	} else {
		if err := db.Where("id = ?", id).First(&topic).Error; err != nil {
			c.JSON(200, gin.H{"msg": []Content{Content{"Id doesn't exist"}}})
		} else {
			if err := c.ShouldBindJSON(&topic); err != nil {
				c.JSON(200, gin.H{"msg": []Content{Content{"Invalid format or datatype"}}})
				c.Abort()
			} else {
				db.Save(&topic)
				c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
			}
		}
	}
}

func GetTopics(c *gin.Context) {
	var topics []Topics
	if err := db.Find(&topics).Error; err != nil {
		c.JSON(200, gin.H{"msg": []Empty{Empty{}}})
	} else {
		for i, topic := range topics {
			var modregs []*ModbusRegisters
			if err = db.Model(&topic).Association("ModRegs").Find(&modregs).Error; err != nil {
				// topics[i].ModRegs = modregs
			} else {
				topics[i].ModRegs = modregs
			}
		}
		c.JSON(200, gin.H{"msg": topics})
	}
}

func DeleteTopics(c *gin.Context) {
	id := c.Params.ByName("id")
	var topic Topics
	if err = db.Where("id = ?", id).First(&topic).Error; err != nil {
		fmt.Println("error...", err)
		c.JSON(200, gin.H{"msg": []Content{Content{"Topic not found"}}})
		c.Abort()
	} else {
		if res := db.Model(&topic).Association("ModRegs").Clear(); res.Error != nil {
			fmt.Println(err)
		} else {
			d := db.Where("id = ?", id).Delete(&topic)
			c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
			_ = d
		}
	}
}

func AssoTopic2Regs(c *gin.Context) {
	var modregids AssoRegsStruct
	id := c.Params.ByName("id")
	if err = c.ShouldBindJSON(&modregids); err != nil {
		c.JSON(200, gin.H{"msg": []Content{Content{"Invalid Data"}}})
	} else {
		var topic Topics
		var modregs []ModbusRegisters
		if err = db.Where("id = ?", id).First(&topic).Error; err != nil {
			fmt.Println("error...", err)
			c.JSON(200, gin.H{"msg": []Content{Content{"Topic not found"}}})
			c.Abort()
		} else {
			if err = db.Where("id in (?)", modregids.ModRegId).Find(&modregs).Error; err != nil {
				c.JSON(200, gin.H{"msg": []Content{Content{"Error in getting Modbus Registers"}}})
				c.Abort()
			} else {
				if len(modregs) > 0 {
					if res := db.Model(&topic).Association("ModRegs").Append(modregs); res.Error != nil {
						fmt.Println(err)
					}
					c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
					c.Abort()
				} else {
					c.JSON(200, gin.H{"msg": []Content{Content{"No Modbus Registers Found"}}})
					c.Abort()
				}
			}
		}
	}
}

func DelTopic2RegsAll(c *gin.Context) {
	id := c.Params.ByName("id")
	var topic Topics
	if err = db.Where("id = ?", id).First(&topic).Error; err != nil {
		fmt.Println("error...", err)
		c.JSON(200, gin.H{"msg": []Content{Content{"Topic not found"}}})
		c.Abort()
	} else {
		if res := db.Model(&topic).Association("ModRegs").Clear(); res.Error != nil {
			fmt.Println(err)
		} else {
			c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
			c.Abort()
		}
	}
}
func DelTopic2Regs(c *gin.Context) {
	var modregids AssoRegsStruct
	id := c.Params.ByName("id")
	if err = c.ShouldBindJSON(&modregids); err != nil {
		c.JSON(200, gin.H{"msg": []Content{Content{"Invalid Data"}}})
	} else {
		var topic Topics
		var modregs []ModbusRegisters
		if err = db.Where("id = ?", id).First(&topic).Error; err != nil {
			fmt.Println("error...", err)
			c.JSON(200, gin.H{"msg": []Content{Content{"Topic not found"}}})
			c.Abort()
		} else {
			if err = db.Where("id in (?)", modregids.ModRegId).Find(&modregs).Error; err != nil {
				c.JSON(200, gin.H{"msg": []Content{Content{"Error in getting Modbus Registers"}}})
				c.Abort()
			} else {
				if len(modregs) > 0 {
					if res := db.Model(&topic).Association("ModRegs").Delete(modregs); res.Error != nil {
						fmt.Println(err)
					}
					c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
					c.Abort()
				} else {
					c.JSON(200, gin.H{"msg": []Content{Content{"No Modbus Registers Found"}}})
					c.Abort()
				}
			}
		}
	}
}

func AssoPort2Regs(c *gin.Context) {
	var modregids AssoRegsStruct
	id := c.Params.ByName("id")
	if err = c.ShouldBindJSON(&modregids); err != nil {
		c.JSON(200, gin.H{"msg": []Content{Content{"Invalid Data"}}})
	} else {
		var port SerialDetails
		var modregs []ModbusRegisters
		if err = db.Where("id = ?", id).First(&port).Error; err != nil {
			fmt.Println("error...", err)
			c.JSON(200, gin.H{"msg": []Content{Content{"Serial Port not found"}}})
			c.Abort()
		} else {
			if err = db.Where("id in (?)", modregids.ModRegId).Find(&modregs).Error; err != nil {
				c.JSON(200, gin.H{"msg": []Content{Content{"Error in getting Modbus Registers"}}})
				c.Abort()
			} else {
				if len(modregs) > 0 {
					if res := db.Model(&port).Association("ModRegs").Append(modregs); res.Error != nil {
						fmt.Println(err)
					}
					c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
					c.Abort()
				} else {
					c.JSON(200, gin.H{"msg": []Content{Content{"No Modbus Registers Found"}}})
					c.Abort()
				}
			}
		}
	}
}

func DelPort2RegsAll(c *gin.Context) {
	id := c.Params.ByName("id")
	var port SerialDetails
	if err = db.Where("id = ?", id).First(&port).Error; err != nil {
		fmt.Println("error...", err)
		c.JSON(200, gin.H{"msg": []Content{Content{"Serial Port not found"}}})
		c.Abort()
	} else {
		if res := db.Model(&port).Association("ModRegs").Clear(); res.Error != nil {
			fmt.Println(err)
		} else {
			c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
			c.Abort()
		}
	}
}
func DelPort2Regs(c *gin.Context) {
	var modregids AssoRegsStruct
	id := c.Params.ByName("id")
	if err = c.ShouldBindJSON(&modregids); err != nil {
		c.JSON(200, gin.H{"msg": []Content{Content{"Invalid Data"}}})
	} else {
		var port SerialDetails
		var modregs []ModbusRegisters
		if err = db.Where("id = ?", id).First(&port).Error; err != nil {
			fmt.Println("error...", err)
			c.JSON(200, gin.H{"msg": []Content{Content{"Serial Port not found"}}})
			c.Abort()
		} else {
			if err = db.Where("id in (?)", modregids.ModRegId).Find(&modregs).Error; err != nil {
				c.JSON(200, gin.H{"msg": []Content{Content{"Error in getting Modbus Registers"}}})
				c.Abort()
			} else {
				if len(modregs) > 0 {
					if res := db.Model(&port).Association("ModRegs").Delete(modregs); res.Error != nil {
						fmt.Println(err)
					}
					c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
					c.Abort()
				} else {
					c.JSON(200, gin.H{"msg": []Content{Content{"No Modbus Registers Found"}}})
					c.Abort()
				}
			}
		}
	}
}

func AssoRegs2Topic(c *gin.Context) {
	var topicids AssoTopicsStruct
	id := c.Params.ByName("id")
	if err = c.ShouldBindJSON(&topicids); err != nil {
		c.JSON(200, gin.H{"msg": []Content{Content{"Invalid Data"}}})
	} else {
		var topics []Topics
		var modreg ModbusRegisters
		if err = db.Where("id = ?", id).First(&modreg).Error; err != nil {
			fmt.Println("error...", err)
			c.JSON(200, gin.H{"msg": []Content{Content{"Modbus Register not found"}}})
			c.Abort()
		} else {
			if err = db.Where("id in (?)", topicids.TopicsId).Find(&topics).Error; err != nil {
				c.JSON(200, gin.H{"msg": []Content{Content{"Error in getting Topics"}}})
				c.Abort()
			} else {
				if len(topics) > 0 {
					if res := db.Model(&modreg).Association("MqTopic").Append(topics); res.Error != nil {
						fmt.Println(err)
					}
					c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
					c.Abort()
				} else {
					c.JSON(200, gin.H{"msg": []Content{Content{"No Topics Found"}}})
					c.Abort()
				}
			}
		}
	}
}

func DelRegs2Topic(c *gin.Context) {
	var topicids AssoTopicsStruct
	id := c.Params.ByName("id")
	if err = c.ShouldBindJSON(&topicids); err != nil {
		c.JSON(200, gin.H{"msg": []Content{Content{"Invalid Data"}}})
	} else {
		var topics []Topics
		var modreg ModbusRegisters
		if err = db.Where("id = ?", id).First(&modreg).Error; err != nil {
			fmt.Println("error...", err)
			c.JSON(200, gin.H{"msg": []Content{Content{"Modbus Register not found"}}})
			c.Abort()
		} else {
			if err = db.Where("id in (?)", topicids.TopicsId).Find(&topics).Error; err != nil {
				c.JSON(200, gin.H{"msg": []Content{Content{"Error in getting Topics"}}})
				c.Abort()
			} else {
				if len(topics) > 0 {
					if res := db.Model(&modreg).Association("MqTopic").Delete(topics); res.Error != nil {
						fmt.Println(err)
					}
					c.JSON(200, gin.H{"msg": []Content{Content{"Done"}}})
					c.Abort()
				} else {
					c.JSON(200, gin.H{"msg": []Content{Content{"No Topics Found"}}})
					c.Abort()
				}
			}
		}
	}
}

// func PostUpdateProfile(c *gin.Context) {
// 	session := sessions.Default(c)
// 	userId := session.Get("user")
// 	if userId == nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Not logged in"})
// 	} else {
// 		var updateUser UserUpdate
// 		c.Bind(&updateUser)

// 		var user User
// 		if err := db.Where("user_name = ?", userId).First(&user).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
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
// 			db.Save(&user)
// 			c.JSON(200, user)
// 		}
// 	}

// }
