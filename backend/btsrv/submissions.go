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

// SubmittedFile contains the DB model for a submitted file and its transcription
type SubmittedFile struct {
	ID             int             `json:"id" db:"id"`
	URL            string          `json:"url" db:"-"`
	Filename       string          `json:"-" db:"filename"`
	Transcriptions []Transcription `json:"transcriptions" db:"-"`
}

// Transcription is the data for a transcribed image
type Transcription struct {
	ID               int       `json:"id" db:"id"`
	FileID           int       `json:"-" db:"submission_file_id"`
	Transcriber      string    `json:"transcriber" db:"transcriber_name"`
	TranscriberEmail string    `json:"transcriber_email" db:"transcriber_email"`
	Text             string    `json:"text" db:"transcription"`
	TranscribedAt    time.Time `json:"transcribed_at" db:"submitted_at"`
	Approved         bool      `json:"approved" db:"approved"`
}

// TableName sets the name of the table in the DB that this struct binds to
func (sub *Transcription) TableName() string {
	return "transcriptions"
}

// SubmissionDetails contains full details of a sibmission
type SubmissionDetails struct {
	Submission
	InstitutionName string          `json:"institution" db:"-"`
	Files           []SubmittedFile `json:"files" db:"-"`
	Tags            []string        `json:"tags" db:"-"`
	NextID          int             `json:"nextId" db:"-"`
	PreviousID      int             `json:"previousId" db:"-"`
}

// ClientSubmission contains the data necessary for a user to make a book traces submission
type ClientSubmission struct {
	Submission
	InstitutionName string   `json:"institution" db:"-" binding:"required"`
	Files           []string `json:"files" binding:"required" db:"-"`
	Tags            []string `json:"tags" db:"-"`
}

// AdminSubmission contains the data necessary for an admin submission update
type AdminSubmission struct {
	Submission
	InstitutionName string   `json:"institution" db:"-" binding:"required"`
	Tags            []string `json:"tags" db:"-"`
}

// Submission is the DB model for a submission. Many of the fields are reused in the ClientSubmission
type Submission struct {
	ID            int       `json:"id"`
	UploadID      string    `json:"uploadId" binding:"required" db:"upload_id"`
	Title         string    `json:"title" binding:"required"`
	Author        string    `json:"author" binding:"required"`
	Publication   string    `json:"publication" db:"publication_info"`
	Library       string    `json:"library"`
	InstitutionID int       `json:"institution_id" db:"institution_id" binding:"required"`
	CallNumber    string    `json:"callNumber" db:"call_number"`
	Description   string    `json:"description"`
	Submitter     string    `json:"submitter" binding:"required" db:"submitter_name"`
	Email         string    `json:"email" binding:"required" db:"submitter_email"`
	SubmittedAt   time.Time `json:"submittedAt" db:"submitted_at"`
	Public        bool      `json:"published" db:"public"`
}

// TableName sets the name of the table in the DB that this struct binds to
func (sub *Submission) TableName() string {
	return "submissions"
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
	detail, err := svc.lookupSubmission(tgtID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, detail)
}

func (svc *ServiceContext) lookupSubmission(tgtID string) (*SubmissionDetails, error) {
	q := svc.DB.NewQuery("select * from submissions where id={:id}")
	q.Bind((dbx.Params{"id": tgtID}))
	sub := &SubmissionDetails{}
	err := q.One(sub)
	if err != nil {
		log.Printf("ERROR: Unable to get details for submission %s:%s", tgtID, err.Error())
		return nil, err
	}

	log.Printf("Get institution for submission %s", tgtID)
	q = svc.DB.NewQuery("select * from institutions where id={:id}")
	q.Bind((dbx.Params{"id": sub.InstitutionID}))
	var inst Institution
	err = q.One(&inst)
	if err != nil {
		log.Printf("WARN: Unable to get institution for %d failed: %s", sub.ID, err.Error())
	} else {
		sub.InstitutionName = inst.Name
	}

	log.Printf("Get next/previous submissions %s", tgtID)
	qs := fmt.Sprintf(`select id from submissions where public = 1 and id > %d limit 1`, sub.ID)
	qn := svc.DB.NewQuery(qs)
	var nextID []int
	err = qn.Column(&nextID)
	if err != nil {
		log.Printf("WARN: Next ID request for %d failed: %s", sub.ID, err.Error())
	} else {
		if len(nextID) > 0 {
			sub.NextID = nextID[0]
		}
	}

	var prevID []int
	qs = fmt.Sprintf(`select id from submissions where public = 1 and id < %d order by id desc limit 1`, sub.ID)
	qp := svc.DB.NewQuery(qs)
	err = qp.Column(&prevID)
	if err != nil {
		log.Printf("WARN: Previous ID request for %d failed: %s", sub.ID, err.Error())
	} else {
		if len(prevID) > 0 {
			sub.PreviousID = prevID[0]
		}
	}

	log.Printf("Get tags for submission %d", sub.ID)
	sub.Tags = make([]string, 0)
	q = svc.DB.NewQuery("select t.name as name from tags t inner join submission_tags st on st.tag_id = t.id where submission_id={:id}")
	q.Bind((dbx.Params{"id": sub.ID}))
	rows, err := q.Rows()
	if err != nil {
		log.Printf("WARNING: Unable to retrieve tags for submission[%d]: %s", sub.ID, err.Error())
	} else {
		for rows.Next() {
			var t struct{ Name string }
			rows.ScanStruct(&t)
			sub.Tags = append(sub.Tags, t.Name)
		}
	}

	log.Printf("Get file URLS for submission %d", sub.ID)
	q = svc.DB.NewQuery("select id,filename from submission_files where submission_id={:id}")
	q.Bind((dbx.Params{"id": sub.ID}))
	rows, err = q.Rows()
	if err != nil {
		log.Printf("WARNING: Unable to retrieve files for submission[%d]: %s", sub.ID, err.Error())
	} else {
		for rows.Next() {
			var f SubmittedFile
			rows.ScanStruct(&f)
			f.URL = fmt.Sprintf("/uploads/%s", f.Filename)

			var transcriptions []Transcription
			tq := svc.DB.NewQuery("select * from transcriptions where submission_file_id={:fid} order by submitted_at desc")
			tq.Bind(dbx.Params{"fid": f.ID})
			terr := tq.All(&transcriptions)
			if terr != nil {
				log.Printf("WARN: unable to get transcriptions for file %d", f.ID)
				f.Transcriptions = make([]Transcription, 0)
			} else {
				f.Transcriptions = transcriptions
			}
			sub.Files = append(sub.Files, f)
		}
	}
	return sub, nil
}

// SubmitTranscription accepts a transcription from the client and stores in the the DB as non-approved
func (svc *ServiceContext) SubmitTranscription(c *gin.Context) {
	var params struct {
		SubmissionID  int    `json:"submissionID"`
		FileID        int    `json:"fileID"`
		Transcription string `json:"transcription"`
		Transcriber   string `json:"name"`
		Email         string `json:"email"`
	}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		log.Printf("ERROR: Invalid transciption submission - %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	q := svc.DB.NewQuery("select * from submissions where id={:id}")
	q.Bind((dbx.Params{"id": params.SubmissionID}))
	var sub SubmissionDetails
	err = q.One(&sub)
	if err != nil {
		log.Printf("ERROR: Unable to get details for submission %d:%s", params.SubmissionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("%s (%s) has submitted a transcription for file %d", params.Email, params.Transcriber, params.FileID)
	model := Transcription{FileID: params.FileID, Transcriber: params.Transcriber,
		TranscriberEmail: params.Email, Text: params.Transcription, TranscribedAt: time.Now()}
	err = svc.DB.Model(&model).Insert()
	if err != nil {
		log.Printf("ERROR: unable to add transcription: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	svc.sendTranscriptionEmail(sub.ID, sub.Title, model)

	c.String(http.StatusOK, "ok")
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
		where public=1 order by id desc limit 5`)
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
			group by s.id order by s.id desc limit %d,%d`, start, pageSize)
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
		url := fmt.Sprintf("/uploads/%s", getThumbFilename(fi.Filename))
		out.Thumbs = append(out.Thumbs, Thumb{SubmissionID: fi.SubID, URL: url})
	}

	c.JSON(http.StatusOK, out)
}

// SubmitForm accepts a user submission
func (svc *ServiceContext) SubmitForm(c *gin.Context) {
	var submission ClientSubmission
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

	log.Printf("Attach files to submission")
	for _, fn := range submission.Files {
		_, err := svc.DB.Insert("submission_files", dbx.Params{
			"submission_id": submission.ID,
			"filename":      fmt.Sprintf("%s/%s/%s", dateDirs, submission.UploadID, fn),
		}).Execute()
		if err != nil {
			log.Printf("WARN: Unable to attach %s to submission %d", fn, submission.ID)
		}
	}

	writeTags(svc.DB, submission.ID, submission.Tags)

	svc.sendSubmissionEmail(submission)

	c.JSON(http.StatusOK, submission)
}

// WriteTags commits submission tags to DB
func writeTags(db *dbx.DB, submissionID int, tags []string) {
	log.Printf("Attach tags to submission")
	for _, tagIDStr := range tags {
		tagID, _ := strconv.Atoi(tagIDStr)
		_, err := db.Insert("submission_tags", dbx.Params{
			"submission_id": submissionID,
			"tag_id":        tagID,
		}).Execute()
		if err != nil {
			log.Printf("WARN: Unable to attach tag %d to submission %d", tagID, submissionID)
		}
	}
}
