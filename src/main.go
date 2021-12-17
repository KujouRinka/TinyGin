package main

import (
	"TinyGin/src/tiny_gin"
	"net/http"
)

func main() {
	r := tiny_gin.New()
	r.GET("/", func(ctx *tiny_gin.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello World</h1>")
	})
	r.GET("/hello", func(ctx *tiny_gin.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})

	r.POST("/login", func(ctx *tiny_gin.Context) {
		ctx.JSON(http.StatusOK, tiny_gin.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})

	r.Run(":8888")
}
