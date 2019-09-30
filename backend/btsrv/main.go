package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Version of the service
const version = "1.1.0"

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
		api.GET("/news", svc.GetNews)
		api.GET("/search", svc.Search)
		api.GET("/thumbs", svc.GetThumbs)
		api.GET("/recents", svc.GetRecents)
		api.GET("/archives", svc.GetArchivesList)
		api.GET("/submissions/:id", svc.GetSubmissionDetail)
		api.POST("/submit", svc.SubmitForm)
		api.POST("/upload", svc.UploadFile)
		api.DELETE("/upload/:file", svc.DeleteUploadedFile)
		admin := api.Group("/admin")
		{
			admin.POST("/news", svc.AuthMiddleware, svc.AddNews)
			admin.DELETE("/news/:id", svc.AuthMiddleware, svc.DeleteNews)
			admin.PUT("/news/:id", svc.AuthMiddleware, svc.UpdateNews)
			admin.POST("/events", svc.AuthMiddleware, svc.AddEvent)
			admin.DELETE("/events/:id", svc.AuthMiddleware, svc.DeleteEvent)
			admin.PUT("/events/:id", svc.AuthMiddleware, svc.UpdateEvent)
			admin.GET("/submissions", svc.AuthMiddleware, svc.GetSubmissions)
			admin.DELETE("/submissions/:id", svc.AuthMiddleware, svc.DeleteSubmission)
			admin.PUT("/submissions/:id/rotate", svc.AuthMiddleware, svc.RotateImage)
			admin.POST("/submissions/:id/publish", svc.AuthMiddleware, svc.PublishSubmission)
			admin.POST("/submissions/:id/unpublish", svc.AuthMiddleware, svc.UnpublishSubmission)
			admin.POST("/submissions/:id", svc.AuthMiddleware, svc.UpdateSubmission)
		}
	}

	// Note: in dev mode, this is never actually used. The front end is served
	// by yarn and it proxies all requests to the API to the routes above
	router.Use(static.Serve("/", static.LocalFile("./public", true)))

	imgDir := fmt.Sprintf("%s/%s", cfg.UploadDir, "submitted")
	log.Printf("Mount %s as /uploads", imgDir)
	router.Use(static.Serve("/uploads", static.LocalFile(imgDir, true)))

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
