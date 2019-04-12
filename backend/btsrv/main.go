package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Version of the service
const version = "1.0.0"

/**
 * MAIN
 */
func main() {
	log.Printf("===> Book Traces Public Submission service staring up <===")

	// Get config params; service port, directories, DB
	cfg := ServiceConfig{}
	cfg.Load()
	svc := ServiceContext{}
	svc.Init(&cfg)

	log.Printf("Setup routes...")
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.Default()
	router.GET("/version", svc.GetVersion)
	router.GET("/healthcheck", svc.HealthCheck)
	router.GET("/authenticate", svc.Authenticate)
	api := router.Group("/api")
	{
		api.GET("/identifier", svc.GetUploadID)
		api.GET("/tags", svc.GetTags)
		api.GET("/events", svc.GetEvents)
		api.POST("/submit", svc.SubmitForm)
		api.POST("/upload", svc.UploadFile)
		api.DELETE("/upload/:file", svc.DeleteUploadedFile)
	}

	// Note: in dev mode, this is never actually used. The front end is served
	// by yarn and it proxies all requests to the API to the routes above
	router.Use(static.Serve("/", static.LocalFile("./public", true)))

	// add a catchall route that renders the index page.
	// based on no-history config setup info here:
	//    https://router.vuejs.org/guide/essentials/history-mode.html#example-server-configurations
	router.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	portStr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Start service v%s on port %s", version, portStr)
	log.Fatal(router.Run(portStr))
}
