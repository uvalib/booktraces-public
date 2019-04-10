package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Tag contains the data for a tag
type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetTags returns a list of tags as JSON
func (svc *ServiceContext) GetTags(c *gin.Context) {
	q := svc.DB.NewQuery("SELECT id, name FROM tags")
	var tags []Tag
	err := q.All(&tags)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrive tags: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, tags)
}
