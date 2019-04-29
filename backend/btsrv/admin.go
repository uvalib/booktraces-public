package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetSubmissions gets one page of admin submissions data
func (svc *ServiceContext) GetSubmissions(c *gin.Context) {
	// Params: page, filter
	pageStr := c.Query("page")
	const pageSize = 50
	page := 1
	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err == nil {
			page = pageInt
		}
	}
	start := (page - 1) * pageSize
	qs := fmt.Sprintf(`select * from submissions order by submitted_at desc limit %d,%d`, start, pageSize)
	q := svc.DB.NewQuery(qs)
	var subs []Submission
	err := q.All(&subs)
	if err != nil {
		log.Printf("ERROR: unable to retrieve submissions: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, subs)
}
