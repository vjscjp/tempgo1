package main

import (
	"crypto/tls"
	"fmt"
	f "github.com/MustWin/gomarathon"
	"github.com/gin-gonic/gin"
	"strings"
	"os"
)

/*var (
	url		   = "https://shipped-tx3-control-01.tx3.shipped-cisco.com/marathon"
	user       = "synthetic-mon"
	pwd        = "VpYdy5abudqkk3Ts"
)*/


type APP struct{
	Appid string
}

type HOST_PORT struct{
	Host string
	Port int
}


func main() {
	os.Setenv("MARATHON_URL","https://shipped-tx3-control-01.tx3.shipped-cisco.com/marathon")
	os.Setenv("MARATHON_USER","synthetic-mon")
	os.Setenv("MARATHON_PASS","VpYdy5abudqkk3Ts")
	
	url := os.Getenv("MARATHON_URL")
	user := os.Getenv("MARATHON_USER")
	pwd := os.Getenv("MARATHON_PASS")
	
	r := gin.Default()
	r.GET("/",func(c *gin.Context){
		c.JSON(200,gin.H{"msg":"ok"})
	})
	auth := f.HttpBasicAuth{user, pwd}
	client, err := f.NewClient(url, &auth, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		fmt.Println("Invalid Details : "+err.Error())
		return
	}
	r.POST("/app",func(c *gin.Context){
		var a APP
		err := c.BindJSON(&a)
		if err != nil{
			c.JSON(400,gin.H{"message":"Bad Request"})
			return
		}
		
		appid := strings.Replace(a.Appid, "/", "", -1)
		resp, err := client.GetApp(appid)
		if err != nil {
			fmt.Println("Invalid Details : "+err.Error())
			return
		}

		c.JSON(200,resp)		
	})
	r.POST("/host_port",func(c *gin.Context){
		var a HOST_PORT
		err := c.BindJSON(&a)
		if err != nil{
			c.JSON(400,gin.H{"message":"Bad Request"})
			return
		}
		
		flag := false
		resp, err := client.ListTasks()
		if err != nil {
			c.JSON(400,gin.H{"message":err})
			return
		}
		for _, t := range resp.Tasks {
			if t.Host == a.Host {
				for _, p := range t.Ports {
					//fmt.Println("DEBUG port ", p)
					if p == a.Port {
						resp, err := client.GetApp(t.AppID)
						if err != nil{
							c.JSON(400,gin.H{"message":"Bad Request"})
							return
						}
						c.JSON(200,resp)		
						flag = true
					}
				}
			}
		}
		if !flag {
			c.JSON(400,gin.H{"message":"No Record Found"})
			return
		}	
	})
	r.Run(":8888")
}
