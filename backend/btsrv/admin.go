package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetSubmissions gets one page of admin submissions data
func (svc *ServiceContext) GetSubmissions(c *gin.Context) {
	c.String(http.StatusOK, "shouldnt get here")
}
