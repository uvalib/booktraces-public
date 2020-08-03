package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// UpdateTranscription updates transcription text
func (svc *ServiceContext) UpdateTranscription(c *gin.Context) {
	var req struct {
		Transcription string
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("ERROR: invalid transcription request - %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	sID := c.Param("id")
	tID := c.Param("tid")
	log.Printf("Update transcription %s for submission %s", tID, sID)
	q := svc.DB.NewQuery("update transcriptions set transcription={:t} where id={:id}")
	q.Bind(dbx.Params{"id": tID})
	q.Bind(dbx.Params{"t": req.Transcription})
	_, err = q.Execute()
	if err != nil {
		log.Printf("ERROR: unable to update transcription %s: %s", tID, err.Error())
	}
	c.String(http.StatusOK, "ok")
}

// DeleteTranscription deletes a transcription from the system
func (svc *ServiceContext) DeleteTranscription(c *gin.Context) {
	sID := c.Param("id")
	tID := c.Param("tid")
	log.Printf("Delete transcription %s from submission %s", tID, sID)

	q := svc.DB.NewQuery("delete from transcriptions where id={:id}")
	q.Bind(dbx.Params{"id": tID})
	_, err := q.Execute()
	if err != nil {
		log.Printf("ERROR: Unabe to delete transcription %s: %s", tID, err.Error())
	}

	sub, err := svc.lookupSubmission(sID)
	if err != nil {
		log.Printf("ERROR: unable to get submission %s after transcription %s approval: %s", sID, tID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, sub)
}

// ApproveTranscription approves a transcription and make it visible to the public
func (svc *ServiceContext) ApproveTranscription(c *gin.Context) {
	sID := c.Param("id")
	tID := c.Param("tid")
	log.Printf("Approve transcription %s from submission %s", tID, sID)

	// Get the submission file that the transcription to be approved is for
	var subID int
	q := svc.DB.NewQuery("select submission_file_id from transcriptions where id={:id}")
	q.Bind(dbx.Params{"id": tID})
	err := q.Row(&subID)
	if err == nil {
		// set any other transcription that points to this file to non-approved. this to ensure only
		// one transcription can be approved at a time
		q = svc.DB.NewQuery("update transcriptions set approved=0 where submission_file_id={:id}")
		q.Bind(dbx.Params{"id": subID})
		q.Execute()
	}

	q = svc.DB.NewQuery("update transcriptions set approved=1 where id={:id}")
	q.Bind(dbx.Params{"id": tID})
	_, err = q.Execute()
	if err != nil {
		log.Printf("ERROR: Unabe to approve transcription %s: %s", tID, err.Error())
	}

	sub, err := svc.lookupSubmission(sID)
	if err != nil {
		log.Printf("ERROR: unable to get submission %s after transcription %s approval: %s", sID, tID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, sub)
}

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

// RotateImage rotatates the image specified by URL on the POST body to teh right by 90 degrees.
// The thumbnail image is also rotated
func (svc *ServiceContext) RotateImage(c *gin.Context) {
	log.Printf("Rotate image requested")
	subID := c.Param("id")
	subPath := fmt.Sprintf("%s/submitted", svc.UploadDir)
	imgPath := strings.Replace(c.Query("url"), "/uploads", subPath, -1)
	log.Printf("Submission %s: rotate image %s", subID, imgPath)
	if _, err := os.Stat(imgPath); err != nil {
		log.Printf("File %s not found", imgPath)
		c.String(http.StatusNotFound, fmt.Sprintf("%s not found", imgPath))
		return
	}

	err := rotateImage(imgPath)
	if err != nil {
		log.Printf("Rotate %s failed: %s", imgPath, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	thumbPath := getThumbFilename(imgPath)
	err = rotateImage(thumbPath)
	if err != nil {
		log.Printf("Rotate THUMB %s failed: %s", thumbPath, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "rotated")
}

func rotateImage(imgPath string) error {
	args := []string{"-rotate", "90", imgPath, imgPath}
	log.Printf("convert args: %+v", args)
	cmd := exec.Command("convert", args...)
	return cmd.Run()
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
	var submission ClientSubmission
	err := c.ShouldBindJSON(&submission)
	if err != nil {
		log.Printf("ERROR: Submission update failed - %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("Received updated submission: %+v", submission)
	err = svc.DB.Model(&submission).Exclude("Files", "Tags", "UploadID", "SubmittedAt", "Public").Update()
	if err != nil {
		log.Printf("ERROR: Unable to update submission - %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Updating all tags")
	dq := svc.DB.NewQuery("delete from submission_tags where submission_id={:id}")
	dq.Bind(dbx.Params{"id": submission.ID})
	dq.Execute()

	writeTags(svc.DB, submission.ID, submission.Tags)

	c.String(http.StatusOK, "ok")
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
