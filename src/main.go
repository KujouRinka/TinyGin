package main

import (
	"TinyGin/src/tiny_gin"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := tiny_gin.New()
	r.Use(tiny_gin.Logger())
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	r.GET("/", func(ctx *tiny_gin.Context) {
		ctx.HTML(http.StatusOK, "css.tmpl", nil)
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
