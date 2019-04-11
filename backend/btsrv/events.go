package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Event contains the data for an event
type Event struct {
	ID          string `json:"id"`
	Date        string `json:"date" db:"event_date"`
	Description string `json:"description"`
}

// GetEvents returns a list of tags as JSON
func (svc *ServiceContext) GetEvents(c *gin.Context) {
	q := svc.DB.NewQuery("SELECT id, event_date, description FROM events")
	var events []Event
	err := q.All(&events)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrive events: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, events)
}
