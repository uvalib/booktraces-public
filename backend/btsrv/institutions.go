package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Institution contains the data for an institution
type Institution struct {
	ID   int    `json:"id"`
	Name string `json:"name" db:"name"  binding:"required"`
}

// TableName sets the name of the table in the DB that this struct binds to
func (e *Institution) TableName() string {
	return "institutions"
}

// GetInstitutions returns a list of tags as JSON
func (svc *ServiceContext) GetInstitutions(c *gin.Context) {
	q := svc.DB.NewQuery("SELECT * FROM institutions")
	var institutions []Institution
	err := q.All(&institutions)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrive institutions: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, institutions)
}

// DeleteInstitution removes the specified institution
func (svc *ServiceContext) DeleteInstitution(c *gin.Context) {
	institutionID := c.Param("id")
	log.Printf("Delete institution %s", institutionID)
	q := svc.DB.NewQuery("delete from institutions where id={:id}")
	q.Bind(dbx.Params{"id": institutionID})
	_, err := q.Execute()
	if err != nil {
		log.Printf("ERROR: unable to delete institution %s:%s", institutionID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

// AddInstitution adds a new instutution
func (svc *ServiceContext) AddInstitution(c *gin.Context) {
	log.Printf("Add new institution")
	var inst Institution
	err := c.ShouldBindJSON(&inst)
	if err != nil {
		log.Printf("ERROR: Unable to parse institution: %s ", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	inst.ID = 0
	err = svc.DB.Model(&inst).Insert()
	if err != nil {
		log.Printf("ERROR: Unable to add institution: %s ", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, inst)
}
