package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/go-sql-driver/mysql"
)

// ServiceContext contains the data
type ServiceContext struct {
	UploadDir     string
	DevAuthUser   string
	DB            *dbx.DB
	BookTracesURL string
	SMTP          SMTPConfig
}

// Init will initialize the service context based on the config parameters
func (svc *ServiceContext) Init(cfg *ServiceConfig) {
	log.Printf("Initializing Service...")
	svc.UploadDir = cfg.UploadDir
	svc.DevAuthUser = cfg.DevAuthUser
	svc.SMTP = cfg.SMTP
	svc.BookTracesURL = cfg.BookTracesURL

	log.Printf("Init DB connection to %s...", cfg.DBHost)
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := dbx.Open("mysql", connectStr)
	if err != nil {
		log.Printf("FATAL: Unable to make connection: %s", err.Error())
		os.Exit(1)
	}
	db.LogFunc = log.Printf
	svc.DB = db
}

// GetVersion reports the version of the serivce
func (svc *ServiceContext) GetVersion(c *gin.Context) {
	c.String(http.StatusOK, "Book Traces Public Submission Service version %s", version)
}

// HealthCheck reports the health of the serivce
func (svc *ServiceContext) HealthCheck(c *gin.Context) {
	q := svc.DB.NewQuery("select id,version from versions order by created_at desc limit 1")
	var version struct {
		ID      string
		Version string
	}
	err := q.One(&version)
	if err != nil {
		log.Printf("ERROR: database check failed: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"alive": "true", "mysql": "false"})
		return
	}
	// gin.H is a shortcut for map[string]interface{}
	c.JSON(http.StatusOK, gin.H{"alive": "true", "mysql": "true"})
}
