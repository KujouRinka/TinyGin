package main

import (
	"TinyGin/src/tiny_gin"
	"net/http"
)

func main() {
	r := tiny_gin.New()
	r.GET("/index", func(ctx *tiny_gin.Context) {
		ctx.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	v1 := r.Group("/v1")
	v1.GET("/", func(ctx *tiny_gin.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello World</h1>")
	})
	v1.GET("/hello", func(ctx *tiny_gin.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})

	v2 := r.Group("/v2")
	v2.GET("/hello/:name", func(ctx *tiny_gin.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Param("name"), ctx.Path)
	})
	v2.POST("/login", func(ctx *tiny_gin.Context) {
		ctx.JSON(http.StatusOK, tiny_gin.H{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})

	r.Run(":8888")
}
