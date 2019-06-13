package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// News contains the data for a news item
type News struct {
	ID        int       `json:"id"`
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
	q := svc.DB.NewQuery(`select * from news order by created_at desc`)
	err := q.All(&data)
	if err != nil {
		log.Printf("ERROR: Unable to get news list: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

// DeleteNews removes the specified news
func (svc *ServiceContext) DeleteNews(c *gin.Context) {
	newsID := c.Param("id")
	log.Printf("Delete news item %s", newsID)
	q := svc.DB.NewQuery("delete from news where id={:id}")
	q.Bind(dbx.Params{"id": newsID})
	_, err := q.Execute()
	if err != nil {
		log.Printf("ERROR: unable to delete news %s:%s", newsID, err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

// UpdateNews updates the title or content of the specified news item
func (svc *ServiceContext) UpdateNews(c *gin.Context) {
	newsIDStr := c.Param("id")
	newsID, _ := strconv.Atoi(newsIDStr)
	log.Printf("Update news item %d", newsID)

	var news News
	err := c.ShouldBindJSON(&news)
	if err != nil {
		log.Printf("ERROR: Unable to parse news: %s ", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	news.ID = newsID
	err = svc.DB.Model(&news).Exclude("ID").Update()
	if err != nil {
		log.Printf("ERROR: Unable to update news: %s ", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}

// AddNews adds a new news
func (svc *ServiceContext) AddNews(c *gin.Context) {
	log.Printf("Add new news")
	var news News
	err := c.BindJSON(&news)
	if err != nil {
		log.Printf("ERROR: Unable to parse news: %s ", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("GOT %+v", news)

	news.ID = 0
	news.CreatedAt = time.Now()
	err = svc.DB.Model(&news).Insert()
	if err != nil {
		log.Printf("ERROR: Unable to add news: %s ", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Added %+v", news)
	c.JSON(http.StatusOK, news)
}
