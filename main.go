package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func getTimeNow() string {
	t := time.Now().UTC().Add(3 * time.Hour)
	f := t.Format("2006-01-02T15:04:05.000")
	return f
}

var x = 1

func main() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		for {
			<-ticker.C
			fmt.Println(x)
			fmt.Println(getTimeNow())
			x++
		}
	}()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{AllowAllOrigins: true}))
	r.GET("/sse", sse)
	r.GET("/get", getRequest)
	r.GET("print", print)
	r.Run(":8080")
}

var users = 0

func sse(ctx *gin.Context) {
	fmt.Println("SSE start")
	users++
	ticker := time.NewTicker(40 * time.Second)
	defer ticker.Stop()
	defer func() { //this never runs and this go routines goes forever
		fmt.Println("User disconnected")
		users--
	}()

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Writer.Flush()
	go func() {
		for {
			<-ticker.C
			fmt.Println("Active Users:", users)
			ctx.SSEvent("ping", "hi")
			ctx.Writer.Flush()
		}
	}()
	<-ctx.Request.Context().Done() //Context.Done() fires a signal when the request is done or gets canceled
}

var k = 1

func print(ctx *gin.Context) {
	fmt.Println(k)
	ctx.Writer.Write([]byte(strconv.Itoa(k)))
	k++
}

func getRequest(ctx *gin.Context) {
	fmt.Println("GET Start")
	<-ctx.Request.Context().Done()
	fmt.Println("GET end") //this never runs too for the same reason
}
