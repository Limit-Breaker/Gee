package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gee"
)

//
func middlewareLog() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Next()
		// Calculate resolution time
		log.Printf("[%d], cost %s in %v for group v1", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	fmt.Println("start Gee ...")

	engine := gee.NewEngine()

	engine.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	engine.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	v1 := engine.Group("/v1")
	// 使用中间件
	v1.Use(middlewareLog())
	{
		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=makabaka
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})

		v1.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/makabaka
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}
	v2 := engine.Group("/v2")
	{
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	engine.Run(":9999")

	fmt.Println("end Gee ...")
}
