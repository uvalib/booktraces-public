package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// GetUploadID will generate an unique token to identify a new submission
// It will be used as a storage subdir for submission files as they are uploaded
func (svc *ServiceContext) GetUploadID(c *gin.Context) {
	id := xid.New()
	c.String(http.StatusOK, id.String())
}

// UploadFile handles raw file uploads from the front end. These will go into a pending directory
func (svc *ServiceContext) UploadFile(c *gin.Context) {
	log.Printf("INFO: checking for upload identifier...")
	uploadID := c.PostForm("uploadID")
	if uploadID == "" {
		log.Printf("ERROR: No upload identifier submitted")
		c.String(http.StatusBadRequest, "upload identifier missing")
		return
	}

	pendingDir := fmt.Sprintf("%s/%s", svc.UploadDir, "pending")
	uploadDir := fmt.Sprintf("%s/%s", pendingDir, uploadID)
	log.Printf("INFO: identifier received. Create upload directory %s", uploadDir)
	err := os.MkdirAll(uploadDir, 0777)
	if err != nil {
		log.Printf("ERROR: unable to create %s: %s", uploadDir, err.Error())
		return
	}

	file, err := c.FormFile("uploadimage")
	if err != nil {
		log.Printf("ERROR: unable to get uploadimage from form: %s", err.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	filename := filepath.Base(file.Filename)
	dest := fmt.Sprintf("%s/%s", uploadDir, filename)
	if _, err := os.Stat(dest); err == nil {
		log.Printf("WARNING: File %s already exists; removing", dest)
		os.Remove(dest)
	}
	log.Printf("INFO: receiving %s", filename)
	if err := c.SaveUploadedFile(file, dest); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}
	log.Printf("INFO: done receiving %s", filename)
	c.String(http.StatusOK, "Submitted")
}

// DeleteUploadedFile will remove a temporary upload file
func (svc *ServiceContext) DeleteUploadedFile(c *gin.Context) {
	tgtFile := c.Param("file")
	uploadID := c.Query("key")
	tgt := fmt.Sprintf("%s/%s/%s/%s", svc.UploadDir, "pending", uploadID, tgtFile)
	log.Printf("INFO: request to delete %s", tgt)
	if _, err := os.Stat(tgt); err == nil {
		delErr := os.Remove(tgt)
		if delErr != nil {
			log.Printf("WARNING: Unable to delete %s: %s", tgt, delErr.Error())
			c.String(http.StatusInternalServerError, delErr.Error())
			return
		}
	} else {
		log.Printf("WARNING: Target file %s does not exist", tgt)
		c.String(http.StatusNotFound, "% not found", tgtFile)
		return
	}
	log.Printf("INFO: deleted %s", tgt)
	c.String(http.StatusOK, "deleted %s", tgtFile)
}
