package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	numbers := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	x := make(map[*string]int)
	go func() {
		for {
			fmt.Println("start len", len(x))
			pref := make([]*string, 0)
			for i := 1; i <= 10; i++ {
				p := &numbers[i-1]
				pref = append(pref, p)
				x[p] = i
			}
			fmt.Println("\n\nshowing\n----------")
			for str, in := range x {
				fmt.Println(*str, in)
			}
			fmt.Println("\n\ndeleting")
			for j := 1; j <= 10; j++ {
				delete(x, pref[j-1])
			}
			fmt.Println("\n\nafter delete show")
			for str, in := range x {
				fmt.Println(*str, in)
			}
			fmt.Println("end len", len(x))
			time.Sleep(10 * time.Second)
		}
	}()
	r := gin.Default()
	r.GET("/normal", normalSSE)
	r.GET("/poly", polySSE)
	r.GET("/get", normalGET)
}
func normalSSE(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Writer.Flush()
	fmt.Println("normal start")
	<-ctx.Request.Context().Done()
	fmt.Println("normal end")
}

func polySSE(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Writer.Flush()
	fmt.Println("poly start")
	<-ctx.Request.Context().Done()
	fmt.Println("poly end")
}
func normalGET(ctx *gin.Context) {
	fmt.Println("normalGET start")
	ctx.Writer.Write([]byte("hii"))
	<-ctx.Request.Context().Done()
	fmt.Println("normalGET end")
}
