package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// News contains the data for a news item
type News struct {
	ID        string    `json:"id"`
	Title     string    `json:"title" db:"title"  binding:"required"`
	Content   string    `json:"content"  binding:"required"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	Published bool      `json:"published" db:"published"`
}

// TableName sets the name of the table in the DB that this struct binds to
func (e *News) TableName() string {
	return "news"
}

// GetNews returns a list of news items as JSON
func (svc *ServiceContext) GetNews(c *gin.Context) {
	var data []News
	q := svc.DB.NewQuery(`select * from news where published=1 order by created_at desc`)
	err := q.All(&data)
	if err != nil {
		log.Printf("ERROR: Unable to get news list: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}
