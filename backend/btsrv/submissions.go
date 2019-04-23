package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Submission contains the data necessary for a user to make a book traces submission
type Submission struct {
	ID          int      `json:"id"`
	UploadID    string   `json:"uploadId" binding:"required" db:"upload_id"`
	Title       string   `json:"title" binding:"required"`
	Author      string   `json:"author" binding:"required"`
	Publication string   `json:"publication" db:"publication_info"`
	Library     string   `json:"library" binding:"required"`
	CallNumber  string   `json:"callNumber" db:"call_number"`
	Description string   `json:"description"`
	Files       []string `json:"files" binding:"required"`
	Submitter   string   `json:"submitter" binding:"required" db:"submitter_name"`
	Email       string   `json:"email" binding:"required" db:"submitter_email"`
	Tags        []string `json:"tags"`
}

// SubmissionThumb contains the ID and URL for a submission thumbnail
type SubmissionThumb struct {
	SubmissionID int    `json:"submissionID"`
	URL          string `json:"url"`
}

// RecentSubmissions contains into from a single page of most recent submissions
type RecentSubmissions struct {
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
	Thumbs   []SubmissionThumb `json:"thumbs"`
}

// TableName sets the name of the table in the DB that this struct binds to
func (c *Submission) TableName() string {
	return "submissions"
}

// GetSubmissions gets one page of recent submissions and returns it along with some context
func (svc *ServiceContext) GetSubmissions(c *gin.Context) {
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
	out := RecentSubmissions{Total: 0, Page: page, PageSize: pageSize}

	log.Printf("Get total submissions")
	tq := svc.DB.NewQuery("select count(*) as total from submissions where approved=1")
	tq.One(&out)

	qs := fmt.Sprintf(`select s.id as sub_id,upload_id,submitted_at,filename from submissions s 
			inner join submission_files f on f.submission_id = s.id where approved = 1
			group by s.id order by submitted_at asc limit %d,%d`, start, pageSize)
	q := svc.DB.NewQuery(qs)
	rows, err := q.Rows()
	if err != nil {
		log.Printf("ERROR: Unable to get recent submissions %s", err.Error())
		c.String(http.StatusInternalServerError, "Unable to retrieve recent submissions")
		return
	}
	for rows.Next() {
		var fi struct {
			SubID     int    `db:"sub_id"`
			UploadID  string `db:"upload_id"`
			Filename  string `db:"filename"`
			Submitted string `db:"submitted_at"`
		}
		rows.ScanStruct(&fi)
		dateStr := strings.Join(strings.Split(fi.Submitted, "-")[:2], "/")
		url := fmt.Sprintf("/uploads/%s/%s/%s", dateStr, fi.UploadID, getThumbFilename(fi.Filename))
		out.Thumbs = append(out.Thumbs, SubmissionThumb{SubmissionID: fi.SubID, URL: url})
	}

	c.JSON(http.StatusOK, out)
}

// SubmitForm accepts a user submission
func (svc *ServiceContext) SubmitForm(c *gin.Context) {
	var submission Submission
	err := c.ShouldBindJSON(&submission)
	if err != nil {
		log.Printf("ERROR: Submission failed - %s", err.Error())
		c.String(http.StatusBadRequest, "All fields are required")
		return
	}
	log.Printf("Received submission: %+v", submission)
	pendingDir := fmt.Sprintf("%s/%s", svc.UploadDir, "pending")
	uploadDir := fmt.Sprintf("%s/%s", pendingDir, submission.UploadID)

	// Submitted dir gets broken up by YYYY/MM before the submissionID
	currTime := time.Now()
	submittedDir := fmt.Sprintf("%s/%s/%s", svc.UploadDir, "submitted", currTime.Format("2006/01"))
	os.MkdirAll(submittedDir, 0777)
	tgtDir := fmt.Sprintf("%s/%s", submittedDir, submission.UploadID)
	log.Printf("Moving pending upload files from %s to %s", uploadDir, tgtDir)
	err = os.Rename(uploadDir, tgtDir)
	if err != nil {
		log.Printf("ERROR: Unable to move pending files to submitted: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	os.Chmod(tgtDir, 0777)

	err = generateThumbnails(tgtDir)
	if err != nil {
		log.Printf("ERROR: Unable to generate thumbnails %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Create DB record of submission")
	err = svc.DB.Model(&submission).Exclude("Files", "Tags").Insert()
	if err != nil {
		log.Printf("ERROR: Unable to add submission rec - %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Attach files to submission")
	for _, fn := range submission.Files {
		_, err := svc.DB.Insert("submission_files", dbx.Params{
			"submission_id": submission.ID,
			"filename":      fn,
		}).Execute()
		if err != nil {
			log.Printf("WARN: Unable to attach %s to submission %d", fn, submission.ID)
		}
	}

	log.Printf("Attach tags to submission")
	for _, tagIDStr := range submission.Tags {
		tagID, _ := strconv.Atoi(tagIDStr)
		_, err := svc.DB.Insert("submission_tags", dbx.Params{
			"submission_id": submission.ID,
			"tag_id":        tagID,
		}).Execute()
		if err != nil {
			log.Printf("WARN: Unable to attach tag %d to submission %d", tagID, submission.ID)
		}
	}

	c.JSON(http.StatusOK, submission)
}

// Generate 150x150x thumbnails for all images (.png and .jpg) in the srcDir
func generateThumbnails(srcDir string) error {
	log.Printf("Generate 150x150 thumbnail for all images in %s", srcDir)
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, imgFile := range files {
		if imgFile.IsDir() {
			continue
		}
		origFn := fmt.Sprintf("%s/%s", srcDir, imgFile.Name())
		os.Chmod(origFn, 0666)
		thumbFn := getThumbFilename(origFn)
		log.Printf("Generate thumbnail file %s...", thumbFn)
		args := []string{"-quiet", "-resize", "150x150^", "-extent", "150x150", "-gravity", "center", origFn, thumbFn}
		log.Printf("convert args: %+v", args)
		cmd := exec.Command("convert", args...)
		err := cmd.Run()
		if err != nil {
			return err
		}
		os.Chmod(thumbFn, 0666)
		log.Printf("Thunbnail %s generated", thumbFn)
	}
	return nil
}

func getThumbFilename(origFn string) string {
	ext := filepath.Ext(origFn)
	return fmt.Sprintf("%s-150x150%s", strings.TrimSuffix(origFn, ext), ext)
}
