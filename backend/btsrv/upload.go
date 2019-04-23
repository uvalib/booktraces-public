package main

import (
	"fmt"
	"io"
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
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/xid"
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

// TableName sets the name of the table in the DB that this struct binds to
func (c *Submission) TableName() string {
	return "submissions"
}

// GetUploadID will generate an unique token to identify a new submission
// It will be used as a storage subdir for submission files as they are uploaded
func (svc *ServiceContext) GetUploadID(c *gin.Context) {
	id := xid.New()
	c.String(http.StatusOK, id.String())
}

// UploadFile handles raw file uploads from the front end. These will go into a pending directory
func (svc *ServiceContext) UploadFile(c *gin.Context) {
	log.Printf("Checking for upload identifier...")
	uploadID := c.PostForm("uploadID")
	if uploadID == "" {
		log.Printf("ERROR: No upload identifier submitted")
		c.String(http.StatusBadRequest, "upload identifier missing")
		return
	}

	log.Printf("Identifier received. Create upload directory.")
	pendingDir := fmt.Sprintf("%s/%s", svc.UploadDir, "pending")
	uploadDir := fmt.Sprintf("%s/%s", pendingDir, uploadID)
	os.MkdirAll(uploadDir, 0777)

	// when chunking is being used, there will be additional form params:
	// dzchunkindex,  dztotalfilesize, dzchunksize,  dztotalchunkcount
	// All sizes are in bytes
	chunkIdx := c.PostForm("dzchunkindex")
	if chunkIdx != "" {
		// this is a chunked file; open in append mode and write to it
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Unable to get form file: %s", err.Error()))
			return
		}
		filename := header.Filename
		log.Printf("Received CHUNKED request to upload %s, chunk %s size %s", filename, chunkIdx, c.PostForm("dzchunksize"))
		dest := fmt.Sprintf("%s/%s", uploadDir, filename)
		if chunkIdx == "0" {
			if _, err := os.Stat(dest); err == nil {
				log.Printf("WARN: File %s already exists; removing", dest)
				os.Remove(dest)
			}
		}
		outFile, err := os.OpenFile(dest, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("unable to receive file %s", err.Error()))
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, file)
	} else {
		// not chunked; just save the file in the temp dir
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		filename := filepath.Base(file.Filename)
		dest := fmt.Sprintf("%s/%s", uploadDir, filename)
		if _, err := os.Stat(dest); err == nil {
			log.Printf("WARN: File %s already exists; removing", dest)
			os.Remove(dest)
		}
		log.Printf("Receiving non-chunked file %s", filename)
		if err := c.SaveUploadedFile(file, dest); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		log.Printf("Done receiving %s", filename)
		c.String(http.StatusOK, "Submitted")
	}
}

// DeleteUploadedFile will remove a temporary upload file
func (svc *ServiceContext) DeleteUploadedFile(c *gin.Context) {
	tgtFile := c.Param("file")
	uploadID := c.Query("key")
	tgt := fmt.Sprintf("%s/%s/%s/%s", svc.UploadDir, "pending", uploadID, tgtFile)
	log.Printf("Request to delete %s", tgt)
	if _, err := os.Stat(tgt); err == nil {
		delErr := os.Remove(tgt)
		if delErr != nil {
			log.Printf("WARN: Unable to delete %s: %s", tgt, delErr.Error())
			c.String(http.StatusInternalServerError, delErr.Error())
			return
		}
	} else {
		log.Printf("WARN: Target file %s does not exist", tgt)
		c.String(http.StatusNotFound, "% not found", tgtFile)
		return
	}
	log.Printf("Deleted %s", tgt)
	c.String(http.StatusOK, "deleted %s", tgtFile)
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
		ext := filepath.Ext(origFn)
		thumbFn := fmt.Sprintf("%s-150x150%s", strings.TrimSuffix(origFn, ext), ext)
		log.Printf("Generate thumbnail file %s...", thumbFn)
		args := []string{"-quiet", "-resize", "150x150^", "-extent", "150x150", "-gravity", "center", origFn, thumbFn}
		log.Printf("convert args: %+v", args)
		cmd := exec.Command("convert", args...)
		err := cmd.Run()
		if err != nil {
			return err
		}
		log.Printf("Thunbnail %s generated", thumbFn)
	}
	return nil
}
