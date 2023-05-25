package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// numbers := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	// x := make(map[*string]int)
	// go func() {
	// 	for {
	// 		fmt.Println("start len", len(x))
	// 		pref := make([]*string, 0)
	// 		for i := 1; i <= 10; i++ {
	// 			p := &numbers[i-1]
	// 			pref = append(pref, p)
	// 			x[p] = i
	// 		}
	// 		fmt.Println("\n\nshowing\n----------")
	// 		for str, in := range x {
	// 			fmt.Println(*str, in)
	// 		}
	// 		fmt.Println("\n\ndeleting")
	// 		for j := 1; j <= 10; j++ {
	// 			delete(x, pref[j-1])
	// 		}
	// 		fmt.Println("\n\nafter delete show")
	// 		for str, in := range x {
	// 			fmt.Println(*str, in)
	// 		}
	// 		fmt.Println("end len", len(x))
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{AllowOrigins: []string{"http://localhost:5173"}}))
	r.GET("/normal", normalSSE)
	r.GET("/poly", polySSE)
	r.GET("/get", normalGET)
	r.Run(":8080")
}
func normalSSE(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Writer.Flush()
	fmt.Println("normal start")
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-time.Tick(10 * time.Second):
				fmt.Println("sending normal")
				ctx.SSEvent("hi", "Normal Hi")
				ctx.Writer.Flush()
				ctx.SSEvent("ping", struct{}{})
				ctx.Writer.Flush()

			case <-ctx.Request.Context().Done():
				done <- struct{}{}
				fmt.Println("normal done case")
				return
			}

		}

	}()
	<-done
}

func polySSE(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Writer.Flush()
	fmt.Println("poly start")
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-time.Tick(10 * time.Second):
				// ctx.SSEvent("hi", "Poly Hi")
				// ctx.Writer.Flush()
				fmt.Println("sending poly")
				ctx.SSEvent("ping", struct{}{})
				ctx.Writer.Flush()
			case <-ctx.Request.Context().Done():
				done <- struct{}{}
				fmt.Println("poly done case")
				return
			}

		}
	}()
	<-done
}
func normalGET(ctx *gin.Context) {
	fmt.Println("normalGET start")
	ctx.Writer.Write([]byte("hii"))
	<-ctx.Request.Context().Done()
	fmt.Println("normalGET end")
}
