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
	ID          int       `json:"id"`
	UploadID    string    `json:"uploadId" binding:"required" db:"upload_id"`
	Title       string    `json:"title" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	Publication string    `json:"publication" db:"publication_info"`
	Library     string    `json:"library" binding:"required"`
	CallNumber  string    `json:"callNumber" db:"call_number"`
	Description string    `json:"description"`
	Files       []string  `json:"files" binding:"required" db:"-"`
	Submitter   string    `json:"submitter" binding:"required" db:"submitter_name"`
	Email       string    `json:"email" binding:"required" db:"submitter_email"`
	Tags        []string  `json:"tags" db:"-"`
	SubmittedAt time.Time `json:"submittedAt" db:"submitted_at"`
	Public      bool      `json:"published" db:"public"`
}

// TableName sets the name of the table in the DB that this struct binds to
func (sub *Submission) TableName() string {
	return "submissions"
}

// WriteFiles commits submission files to DB
func (sub *Submission) WriteFiles(db *dbx.DB, dateDirs string) {
	log.Printf("Attach files to submission")
	for _, fn := range sub.Files {
		_, err := db.Insert("submission_files", dbx.Params{
			"submission_id": sub.ID,
			"filename":      fmt.Sprintf("%s/%s/%s", dateDirs, sub.UploadID, fn),
		}).Execute()
		if err != nil {
			log.Printf("WARN: Unable to attach %s to submission %d", fn, sub.ID)
		}
	}
}

// WriteTags commits submission tags to DB
func (sub *Submission) WriteTags(db *dbx.DB) {
	log.Printf("Attach tags to submission")
	for _, tagIDStr := range sub.Tags {
		tagID, _ := strconv.Atoi(tagIDStr)
		_, err := db.Insert("submission_tags", dbx.Params{
			"submission_id": sub.ID,
			"tag_id":        tagID,
		}).Execute()
		if err != nil {
			log.Printf("WARN: Unable to attach tag %d to submission %d", tagID, sub.ID)
		}
	}
}

// GetTags retrieves the list of tags associated with this submission
func (sub *Submission) GetTags(db *dbx.DB) {
	q := db.NewQuery("select t.name as name from tags t inner join submission_tags st on st.tag_id = t.id where submission_id={:id}")
	q.Bind((dbx.Params{"id": sub.ID}))
	rows, err := q.Rows()
	if err != nil {
		log.Printf("WARNING: Unable to retrieve tags for submission[%d]: %s", sub.ID, err.Error())
		return
	}
	for rows.Next() {
		var t struct{ Name string }
		rows.ScanStruct(&t)
		sub.Tags = append(sub.Tags, t.Name)
	}
}

// GetFileURLs retrieves the list of URLS for files associated with this submission
func (sub *Submission) GetFileURLs(db *dbx.DB) {
	log.Printf("Getting file URLS for submission %d", sub.ID)
	q := db.NewQuery("select filename from submission_files where submission_id={:id}")
	q.Bind((dbx.Params{"id": sub.ID}))
	rows, err := q.Rows()
	if err != nil {
		log.Printf("WARNING: Unable to retrieve files for submission[%d]: %s", sub.ID, err.Error())
		return
	}

	for rows.Next() {
		var f struct{ Filename string }
		rows.ScanStruct(&f)
		sub.Files = append(sub.Files, fmt.Sprintf("/uploads/%s", f.Filename))
	}
}

// GetArchivesList will get a sorted list of archives dates
func (svc *ServiceContext) GetArchivesList(c *gin.Context) {
	type Archives struct {
		Display  string `json:"displayDate"`
		Internal string `json:"internalDate"`
	}
	var data []Archives
	q := svc.DB.NewQuery(`select distinct DATE_FORMAT(submitted_at,'%M %Y') display, DATE_FORMAT(submitted_at,'%Y-%m') as internal
		 from submissions where public=1 order by submitted_at desc`)
	err := q.All(&data)
	if err != nil {
		log.Printf("ERROR: Unable to get archives list: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetSubmissionDetail gets full details of a submission
func (svc *ServiceContext) GetSubmissionDetail(c *gin.Context) {
	tgtID := c.Param("id")
	q := svc.DB.NewQuery("select * from submissions where id={:id}")
	q.Bind((dbx.Params{"id": tgtID}))
	var sub Submission
	err := q.One(&sub)
	if err != nil {
		log.Printf("ERROR: Unable to get details for submission %s:%s", tgtID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	sub.GetTags(svc.DB)
	sub.GetFileURLs(svc.DB)
	c.JSON(http.StatusOK, sub)
}

// GetRecents gets some details on the 5 most recent submissions
func (svc *ServiceContext) GetRecents(c *gin.Context) {
	type Recent struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		SubmittedAt string `json:"submittedAt" db:"submitted"`
	}
	var data []Recent
	q := svc.DB.NewQuery(`select id,title,DATE_FORMAT(submitted_at,'%M %d, %Y') submitted from submissions 
		where public=1 order by submitted_at desc limit 5`)
	err := q.All(&data)
	if err != nil {
		log.Printf("ERROR: Unable to get archives list: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetThumbs gets one page of recent submissions and returns it along with some context
func (svc *ServiceContext) GetThumbs(c *gin.Context) {
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

	type Thumb struct {
		SubmissionID int    `json:"submissionID"`
		URL          string `json:"url"`
	}
	type RecentThumbs struct {
		Total    int     `json:"total"`
		Page     int     `json:"page"`
		PageSize int     `json:"pageSize"`
		Thumbs   []Thumb `json:"thumbs"`
	}

	out := RecentThumbs{Total: 0, Page: page, PageSize: pageSize}

	log.Printf("Get total submissions")
	tq := svc.DB.NewQuery("select count(*) as total from submissions where public=1")
	tq.One(&out)

	qs := fmt.Sprintf(`select s.id as sub_id,upload_id,submitted_at,filename from submissions s 
			inner join submission_files f on f.submission_id = s.id where public = 1
			group by s.id order by submitted_at desc limit %d,%d`, start, pageSize)
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
		url := fmt.Sprintf("/uploads/%s", GetThumbFilename(fi.Filename))
		out.Thumbs = append(out.Thumbs, Thumb{SubmissionID: fi.SubID, URL: url})
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
	dateDirs := currTime.Format("2006/01")
	submittedDir := fmt.Sprintf("%s/%s/%s", svc.UploadDir, "submitted", dateDirs)
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
	submission.SubmittedAt = time.Now()
	err = svc.DB.Model(&submission).Exclude("Files", "Tags").Insert()
	if err != nil {
		log.Printf("ERROR: Unable to add submission rec - %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	submission.WriteFiles(svc.DB, dateDirs)
	submission.WriteTags(svc.DB)

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
		thumbFn := GetThumbFilename(origFn)
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

// GetThumbFilename takes a filename or path to a file and converts it to a 150x150 thumb
func GetThumbFilename(origFn string) string {
	ext := filepath.Ext(origFn)
	// NOTE: the old WP site makes all derivatives extensions lowercase regardless of original.
	// Make it the same here. This also controls the naming of the thumbs generated by this service.
	return fmt.Sprintf("%s-150x150%s", strings.TrimSuffix(origFn, ext), strings.ToLower(ext))
}
