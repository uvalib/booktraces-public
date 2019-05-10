package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	log.Printf("Publishing submission %s", subID)
	_, err := svc.DB.Update("submissions",
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

// UnpublishSubmission makes a previoulsy public submission private
func (svc *ServiceContext) UnpublishSubmission(c *gin.Context) {
	subID := c.Param("id")
	log.Printf("UNPublishing submission %s", subID)
	_, err := svc.DB.Update("submissions",
		dbx.Params{"public": 0},
		dbx.HashExp{"id": subID}).Execute()
	if err != nil {
		log.Printf("ERROR: Unable to unpublish submission %s: %s", subID, err.Error())
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
		Tags        *string   `json:"tags" db:"tags"`
		SubmittedAt time.Time `json:"submittedAt" db:"submitted_at"`
		Published   int       `json:"published" db:"public"`
	}
	type SubmissionsPage struct {
		Total         int             `json:"total"`
		FilteredTotal int             `json:"filteredTotal"`
		Page          int             `json:"page"`
		PageSize      int             `json:"pageSize"`
		Submissions   []SubmissionRow `json:"submissions"`
	}
	out := SubmissionsPage{Total: 0, Page: page, PageSize: pageSize}

	log.Printf("Get total submissions")
	tq := svc.DB.NewQuery("select count(*) as total from submissions")
	tq.Row(&out.Total)

	// Generate the basic SQL strings used to get submission rows
	selSQL := "select s.id as id, title, author, submitted_at, public, group_concat(t.name) as tags"
	fromSQL := ` from submissions s
		left outer join submission_tags st on st.submission_id = s.id
		left outer join tags t on t.id = st.tag_id`
	qs := selSQL + fromSQL
	groupQS := " group by s.id"
	pageQS := fmt.Sprintf(" order by submitted_at desc limit %d,%d", start, pageSize)

	// Check for and apply and filter / query params
	qParam := strings.TrimSpace(c.Query("q"))
	qQuery := ""
	if qParam != "" {
		log.Printf("Filter submission by query string [%s]", qParam)
		qParam = "%" + qParam + "%"
		qQuery = ` where (t.name like {:q} or s.title like {:q} or s.author like {:q}
			or s.publication_info like {:q} or s.library like {:q} or s.call_number like {:q}
			or s.description like {:q} or s.description like {:q} or s.submitted_at like {:q})`
		qs += qQuery
	}

	tParam := strings.TrimSpace(c.Query("t"))
	if tParam != "" {
		log.Printf("Filter submission by tag [%s]", tParam)
		// To ensure all tags are included in result, can't use where clause.
		// If used, only the single matching tag is returned. Instead add a having
		// clause to the group by. The is leaves all tags in the results and matches
		// on the CSV tag list instead.
		groupQS += " having Find_In_Set({:t}, tags)"
	}

	if qParam != "" || tParam != "" {
		countQS := "select count(distinct s.id) as filtered_cnt " + fromSQL
		if qQuery != "" {
			countQS += qQuery
		}

		// Since all of the tags are not required for a simple match count,
		// the weird group by and having find_in_set is not needed.
		// Just a simple where will work.
		if tParam != "" {
			if strings.Contains(countQS, " where ") {
				countQS += " and "
			} else {
				countQS += " where "
			}
			countQS += " t.name={:t}"
		}

		log.Printf("Get filtered total")
		cq := svc.DB.NewQuery(countQS)
		cq.Bind(dbx.Params{"q": qParam, "t": tParam})
		cq.Row(&out.FilteredTotal)
	}

	log.Printf("Get one page of submission data")
	qs = qs + groupQS + pageQS
	q := svc.DB.NewQuery(qs)
	q.Bind(dbx.Params{"q": qParam, "t": tParam})

	err := q.All(&out.Submissions)
	if err != nil {
		log.Printf("ERROR: Unable to get submisions: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, out)
}
