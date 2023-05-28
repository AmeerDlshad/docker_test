package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func getTimeNow() string {
	t := time.Now().UTC().Add(3 * time.Hour)
	f := t.Format("2006-01-02T15:04:05.000")
	return f
}

// `-F "file=@Dockerfile"`, `https://api.anonfiles.com/upload?token=e6cf3cde4b89f244`
func main() {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Exec("create table if not exists users(name text)")
	// ticker := time.NewTicker(10 * time.Minute)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{AllowAllOrigins: true}))
	r.GET("/users", getUsers)
	r.GET("/add", addUser)
	r.GET("/sse", sse)
	r.GET("/get", getRequest)
	r.GET("print", print)
	r.Run(":8080")
}

var users = 0

func getUsers(ctx *gin.Context) {
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	rows, _ := db.Query("select * from users")
	u := make([]any, 0)
	for rows.Next() {
		var x string
		rows.Scan(&x)
		u = append(u, x)
	}

	ctx.JSON(200, u)
}
func addUser(ctx *gin.Context) {
	name := ctx.Query("u")
	fmt.Println(name)
	db, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = db.Exec("insert into users values(?)", name)
	if err != nil {
		fmt.Println(err)
	}
}
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
