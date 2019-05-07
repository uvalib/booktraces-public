package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// DeleteSubmission deletes all aspected of a submission from the system
func (svc *ServiceContext) DeleteSubmission(c *gin.Context) {
	subID := c.Param("id")
	log.Printf("Delete submission %s", subID)
	q := svc.DB.NewQuery("select upload_id, submitted_at from submissions where id={:id}")
	q.Bind(dbx.Params{"id": subID})
	var sub struct {
		UploadID    string    `db:"upload_id"`
		SubmittedAt time.Time `db:"submitted_at"`
	}
	err := q.One(&sub)
	if err != nil {
		log.Printf("ERROR: unable to find submission %s:%s", subID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	tgtDir := fmt.Sprintf("%s/%s/%s/%s", svc.UploadDir, "submitted", sub.SubmittedAt.Format("2006/01"), sub.UploadID)
	log.Printf("Submission %s found; removing submitted images from %s", subID, tgtDir)
	err = os.RemoveAll(tgtDir)
	if err != nil {
		log.Printf("WARN: Could not remove submitted files from %s", tgtDir)
	}

	log.Printf("Removing submission %s from DB", subID)
	_, err = svc.DB.Delete("submissions", dbx.HashExp{"id": subID}).Execute()
	if err != nil {
		log.Printf("ERROR: Unable to delete submission %s:%s", subID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "deleted")
}

// PublishSubmission publishes a previously unpublished submission
func (svc *ServiceContext) PublishSubmission(c *gin.Context) {
	subID := c.Param("id")
	var posted struct {
		UserID string `json:"userID" binding:"required"`
	}
	err := c.BindJSON(&posted)
	if err != nil {
		log.Printf("ERROR: UserID not found in publish")
		c.String(http.StatusBadRequest, "UserID is required")
		return
	}
	log.Printf("User %s is publishing submission %s", posted.UserID, subID)
	_, err = svc.DB.Update("submissions",
		dbx.Params{"public": 1},
		dbx.HashExp{"id": subID}).Execute()
	if err != nil {
		log.Printf("ERROR: Unable to publish submission %s: %s", subID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Submission %s has been published", subID)
	c.String(http.StatusOK, "ok")
}

// UpdateSubmission updates the details (aside from publish status) of a submission
func (svc *ServiceContext) UpdateSubmission(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not yet")
}

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
		Tags        string    `json:"tags" db:"tags"`
		SubmittedAt time.Time `json:"submittedAt" db:"submitted_at"`
		Published   int       `json:"published" db:"public"`
	}
	type SubmissionsPage struct {
		Total       int             `json:"total"`
		Page        int             `json:"page"`
		PageSize    int             `json:"pageSize"`
		Submissions []SubmissionRow `json:"submissions"`
	}
	out := SubmissionsPage{Total: 0, Page: page, PageSize: pageSize}

	log.Printf("Get total submissions")
	tq := svc.DB.NewQuery("select count(*) as total from submissions where public=1")
	tq.One(&out)

	qs := fmt.Sprintf(`select s.id as id, title, author, submitted_at, public, group_concat(t.name) as tags
		from submissions s, submission_tags st, tags t
		where st.submission_id=s.id and t.id=st.tag_id group by s.id 
		order by submitted_at desc limit %d,%d`, start, pageSize)
	q := svc.DB.NewQuery(qs)
	err := q.All(&out.Submissions)
	if err != nil {
		log.Printf("ERROR: Unable to get submisions: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, out)
}
