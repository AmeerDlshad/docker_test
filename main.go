package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
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

var k = func() int {
	x, _ := os.ReadFile("f.txt")
	y := string(x)
	a, _ := strconv.Atoi(y)

	return a
}()

func print(ctx *gin.Context) {
	fmt.Println(k)
	ctx.Writer.Write([]byte(strconv.Itoa(k)))
	k++
	os.WriteFile("f.txt", []byte(strconv.Itoa(k)), os.ModePerm)
}

func getRequest(ctx *gin.Context) {
	fmt.Println("GET Start")
	<-ctx.Request.Context().Done()
	fmt.Println("GET end") //this never runs too for the same reason
}
