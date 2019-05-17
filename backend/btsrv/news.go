package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

// News contains the data for a news item
type News struct {
	ID        string    `json:"id"`
	Title     string    `json:"title" db:"title"  binding:"required"`
	Content   string    `json:"content"  binding:"required"`
	CreatedAt time.Time `json:"createdAt" db:"created_ad"`
}

// TableName sets the name of the table in the DB that this struct binds to
func (e *News) TableName() string {
	return "news"
}

// GetNews returns a list of news items as JSON
func (svc *ServiceContext) GetNews(c *gin.Context) {
}
