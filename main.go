package main

import (
	"github.com/fitnis/api-gateway/proxy"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	proxy.RegisterRoutes(router)
	router.Run(":8080") // Public API port
}
