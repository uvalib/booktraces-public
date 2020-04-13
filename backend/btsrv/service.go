package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/go-sql-driver/mysql"
)

// ServiceContext contains the data
type ServiceContext struct {
	UploadDir     string
	DevAuthUser   string
	DB            *dbx.DB
	BookTracesURL string
	SMTP          SMTPConfig
}

// Init will initialize the service context based on the config parameters
func (svc *ServiceContext) Init(cfg *ServiceConfig) {
	log.Printf("Initializing Service...")
	svc.UploadDir = cfg.UploadDir
	svc.DevAuthUser = cfg.DevAuthUser
	svc.SMTP = cfg.SMTP
	svc.BookTracesURL = cfg.BookTracesURL

	log.Printf("Init DB connection to %s...", cfg.DBHost)
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := dbx.Open("mysql", connectStr)
	if err != nil {
		log.Printf("FATAL: Unable to make connection: %s", err.Error())
		os.Exit(1)
	}
	db.LogFunc = log.Printf
	svc.DB = db
}

// GetVersion reports the version of the serivce
func (svc *ServiceContext) GetVersion(c *gin.Context) {
	c.String(http.StatusOK, "Book Traces Public Submission Service version %s", version)
}

// HealthCheck reports the health of the serivce
func (svc *ServiceContext) HealthCheck(c *gin.Context) {
	q := svc.DB.NewQuery("select id,version from versions order by created_at desc limit 1")
	var version struct {
		ID      string
		Version string
	}
	err := q.One(&version)
	if err != nil {
		log.Printf("ERROR: database check failed: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"alive": "true", "mysql": "false"})
		return
	}
	// gin.H is a shortcut for map[string]interface{}
	c.JSON(http.StatusOK, gin.H{"alive": "true", "mysql": "true"})
}

func (svc *ServiceContext) sendTranscriptionEmail(subID int, title string, transcription Transcription) {
	log.Printf("Sending transcription notification email for transcription %d", transcription.ID)
	var vars struct {
		BookTracesURL string
		SubmissionID  int
		Title         string
		Transcription Transcription
	}
	vars.BookTracesURL = svc.BookTracesURL
	vars.SubmissionID = subID
	vars.Transcription = transcription
	vars.Title = title
	var renderedEmail bytes.Buffer
	tpl := template.Must(template.New("transcription.txt").ParseFiles("templates/transcription.txt"))
	err := tpl.Execute(&renderedEmail, vars)
	if err != nil {
		log.Printf("ERROR: Unable to render transcription noty email: %s", err.Error())
		return
	}
	svc.sendEmail("Subject: Book Traces Transcription\n", renderedEmail.String())
}

func (svc *ServiceContext) sendSubmissionEmail(submission ClientSubmission) {
	log.Printf("Sending submission notification email for submission %d", submission.ID)
	var vars struct {
		BookTracesURL string
		Submission    ClientSubmission
	}
	vars.BookTracesURL = svc.BookTracesURL
	vars.Submission = submission
	var renderedEmail bytes.Buffer
	tpl := template.Must(template.New("submission.txt").ParseFiles("templates/submission.txt"))
	err := tpl.Execute(&renderedEmail, vars)
	if err != nil {
		log.Printf("ERROR: Unable to render submission noty email: %s", err.Error())
		return
	}
	svc.sendEmail("Subject: Book Traces Submisssion\n", renderedEmail.String())
}

func (svc *ServiceContext) sendEmail(subject string, emailBody string) {
	log.Printf("Generate SMTP message")
	// Per: https://stackoverflow.com/questions/36485857/sending-emails-with-name-email-from-go
	// sending addresses like 'user name <email.com>' does not work with the default
	// mail package. Leaving at just email address for now. Can revisit after meetings
	/// about functionality.
	// toAddr := mail.Address{Name: emailMap[reserveReq.Request.Library], Address: svc.CourseReserveEmail}
	to := []string{svc.SMTP.To}
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	toHdr := fmt.Sprintf("To: %s\n", strings.Join(to, ","))
	msg := []byte(subject + toHdr + mime + emailBody)

	if svc.SMTP.DevMode {
		log.Printf("Email is in dev mode. Logging message instead of sending")
		log.Printf("==================================================")
		log.Printf("%s", msg)
		log.Printf("==================================================")
	} else {
		log.Printf("Sending noty email to %s", strings.Join(to, ","))
		err := smtp.SendMail(fmt.Sprintf("%s:%d", svc.SMTP.Host, svc.SMTP.Port), nil, "no-reply@virginia.edu", to, msg)
		if err != nil {
			log.Printf("ERROR: Unable to send receipt email: %s", err.Error())
		}
	}
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

// getThumbFilename takes a filename or path to a file and converts it to a 150x150 thumb
func getThumbFilename(origFn string) string {
	ext := filepath.Ext(origFn)
	// NOTE: the old WP site makes all derivatives extensions lowercase regardless of original.
	// Make it the same here. This also controls the naming of the thumbs generated by this service.
	return fmt.Sprintf("%s-150x150%s", strings.TrimSuffix(origFn, ext), strings.ToLower(ext))
}
