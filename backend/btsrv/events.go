package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Event contains the data for an event
type Event struct {
	ID          string `json:"id"`
	Date        string `json:"date" db:"event_date"  binding:"required"`
	Description string `json:"description"  binding:"required"`
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

// DeleteEvent removes the specified event
func (svc *ServiceContext) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")
	log.Printf("Delete event %s", eventID)
	q := svc.DB.NewQuery("delete from events where id={:id}")
	q.Bind(dbx.Params{"id": eventID})
	_, err := q.Execute()
	if err != nil {
		log.Printf("ERROR: unable to delete event %s:%s", eventID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

// UpdateEvent updates the date or description of the specified event
func (svc *ServiceContext) UpdateEvent(c *gin.Context) {
	eventID := c.Param("id")
	log.Printf("Update event %s", eventID)

	var event Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		log.Printf("ERROR: Unable to parse event: %s ", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	event.ID = eventID
	err = svc.DB.Model(&event).Exclude("ID").Update()
	if err != nil {
		log.Printf("ERROR: Unable to update event: %s ", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

// AddEvent adds a new event
func (svc *ServiceContext) AddEvent(c *gin.Context) {
	log.Printf("Add new event")
	var event Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		log.Printf("ERROR: Unable to parse event: %s ", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = svc.DB.Model(&event).Exclude("ID").Insert()
	if err != nil {
		log.Printf("ERROR: Unable to add event: %s ", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, event)
}
