package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
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

	type SubmissionRow struct {
		ID          int       `json:"id" db:"id"`
		Title       string    `json:"title" db:"title"`
		Author      string    `json:"author" db:"author"`
		Tags        []string  `json:"tags" db:"-"`
		SubmittedAt time.Time `json:"submittedAt" db:"submitted_at"`
		Published   bool      `json:"published" db:"-"`
	}
	type SubmissionsPage struct {
		Total       int             `json:"total"`
		Page        int             `json:"page"`
		PageSize    int             `json:"pageSize"`
		Submissions []SubmissionRow `json:"submissions"`
	}
	out := SubmissionsPage{Total: 0, Page: page, PageSize: pageSize}

	log.Printf("Get total submissions")
	tq := svc.DB.NewQuery("select count(*) as total from submissions where approved=1")
	tq.One(&out)

	qs := fmt.Sprintf(`select s.id as id, title, author, submitted_at, approved, t.name as tag 
		from submissions s, submission_tags st, tags t
		where st.submission_id=s.id and t.id=st.tag_id  
		order by submitted_at desc limit %d,%d`, start, pageSize)
	q := svc.DB.NewQuery(qs)
	rows, err := q.Rows()
	if err != nil {
		log.Printf("ERROR: Unable to get submissions %s", err.Error())
		c.String(http.StatusInternalServerError, "Unable to retrieve submissions")
		return
	}
	for rows.Next() {
		var sr SubmissionRow
		rowErr := rows.ScanStruct(&sr)
		if rowErr != nil {
			log.Printf("ERROR: Unable to get submissions row %s", rowErr.Error())
			c.String(http.StatusInternalServerError, "Unable to retrieve submissions")
			return
		}
		out.Submissions = append(out.Submissions, sr)
	}

	c.JSON(http.StatusOK, out)
}
