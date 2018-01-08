package main

import "github.com/gin-gonic/gin"

func main() {
	gin.SetMode("release")
	r := newRouter(newLoader())
	r.Use(gin.Recovery())
	r.Run(":8080")
}
