package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Pedagogy contains the data for a news item
type Pedagogy struct {
	ID        int        `json:"id"`
	Key       string     `json:"key" db:"access_key"  binding:"required"`
	Title     string     `json:"title" db:"title"  binding:"required"`
	Content   string     `json:"content"  binding:"required"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

// GetPedagogyDocuments gets a list of pedagogy documents for the admin interface
func (svc *ServiceContext) GetPedagogyDocuments(c *gin.Context) {
	log.Printf("Get all pedagogy documents")
	var docs []Pedagogy
	q := svc.DB.NewQuery(`select * from pedagogy order by title asc`)
	err := q.All(&docs)
	if err != nil {
		log.Printf("ERROR: Unable to get pedagogy documents: %s", err.Error())
		c.JSON(http.StatusNotFound, "no documents found")
		return
	}
	c.JSON(http.StatusOK, docs)
}

// GetPedagogy gets pedagogy content by access key and returns it as a JSON object along with a title
func (svc *ServiceContext) GetPedagogy(c *gin.Context) {
	pedagogyKey := c.Param("id")
	var item Pedagogy
	q := svc.DB.NewQuery(`select * from pedagogy where access_key={:key}`)
	q.Bind(dbx.Params{"key": pedagogyKey})
	err := q.One(&item)
	if err != nil {
		log.Printf("ERROR: Unable to get pedagogy %s: %s", pedagogyKey, err.Error())
		c.JSON(http.StatusNotFound, fmt.Sprintf("%s not found", pedagogyKey))
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeletePedagogy will remove a pedagogy document from the system - except the protected INDEX
func (svc *ServiceContext) DeletePedagogy(c *gin.Context) {
	key := c.Param("id")
	log.Printf("Delete pedagogy  %s", key)
	if key == "index" {
		log.Printf("ERROR: unable to index pedagogy")
		c.String(http.StatusBadRequest, "cannot delete index")
		return
	}

	q := svc.DB.NewQuery("delete from pedagogy where access_key={:key}")
	q.Bind(dbx.Params{"key": key})
	_, err := q.Execute()
	if err != nil {
		log.Printf("ERROR: unable to delete %s:%s", key, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

// UpdatePedagogy will update the title/content of a pedagogy document
func (svc *ServiceContext) UpdatePedagogy(c *gin.Context) {
	key := c.Param("id")
	log.Printf("Update pedagogy %s", key)

	var doc Pedagogy
	err := c.ShouldBindJSON(&doc)
	if err != nil {
		log.Printf("ERROR: Unable to parse pedagogy document: %s ", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	now := time.Now()
	doc.UpdatedAt = &now
	err = svc.DB.Model(&doc).Exclude("CreatedAt").Update()
	if err != nil {
		log.Printf("ERROR: Unable to update pedagogy document %s: %s ", key, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, doc)
}

// AddPedagogy creates a new pedagogy document
func (svc *ServiceContext) AddPedagogy(c *gin.Context) {
	log.Printf("Add new pedagogy document")
	var doc Pedagogy
	err := c.BindJSON(&doc)
	if err != nil {
		log.Printf("ERROR: Unable to parse new pedagogy document: %s ", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	doc.ID = 0
	doc.CreatedAt = time.Now()
	err = svc.DB.Model(&doc).Insert()
	if err != nil {
		log.Printf("ERROR: Unable to add pedagogy document: %s ", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, doc)
}
