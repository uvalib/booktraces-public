package main

import (
	"flag"
	"log"
	"os"
)

// SMTPConfig wraps up all of the smpt configuration
type SMTPConfig struct {
	Host    string
	Port    int
	To      string
	DevMode bool
}

// ServiceConfig defines all of the archives transfer service configuration paramaters
type ServiceConfig struct {
	DBHost        string
	DBName        string
	DBUser        string
	DBPass        string
	DBPort        int
	Port          int
	UploadDir     string
	DevAuthUser   string
	BookTracesURL string
	SMTP          SMTPConfig
}

// Load will load the service configuration from env/cmdline
func (cfg *ServiceConfig) Load() {
	log.Printf("INFO: loading configuration...")

	flag.StringVar(&cfg.DBHost, "dbhost", "", "DB Host (required)")
	flag.StringVar(&cfg.DBName, "dbname", "", "DB Name (required)")
	flag.StringVar(&cfg.DBUser, "dbuser", "", "DB User (required)")
	flag.StringVar(&cfg.DBPass, "dbpass", "", "DB Password (required)")
	flag.IntVar(&cfg.DBPort, "dbport", 3306, "DB Port")

	flag.StringVar(&cfg.DevAuthUser, "devuser", "", "Authorized computing id for dev")
	flag.IntVar(&cfg.Port, "port", 8080, "Service port (default 8080)")
	flag.StringVar(&cfg.UploadDir, "upload", "./submissions", "Upload directory")
	flag.StringVar(&cfg.BookTracesURL, "url", "https://booktraces.org", "Booktraces URL")

	// SMTP settings
	flag.StringVar(&cfg.SMTP.Host, "smtphost", "", "SMTP Host")
	flag.IntVar(&cfg.SMTP.Port, "smtpport", 0, "SMTP Port")
	flag.StringVar(&cfg.SMTP.To, "smtpto", "", "email to receive submission emails")
	flag.BoolVar(&cfg.SMTP.DevMode, "stubsmtp", false, "Log email insted of sending (dev mode)")

	flag.Parse()

	log.Printf("[CONFIG] port          = [%d]", cfg.Port)
	log.Printf("[CONFIG] dbhost        = [%s]", cfg.DBHost)
	log.Printf("[CONFIG] dbport        = [%d]", cfg.DBPort)
	log.Printf("[CONFIG] dbname        = [%s]", cfg.DBName)
	log.Printf("[CONFIG] dbuser        = [%s]", cfg.DBUser)
	log.Printf("[CONFIG] upload        = [%s]", cfg.UploadDir)
	log.Printf("[CONFIG] url           = [%s]", cfg.BookTracesURL)
	log.Printf("[CONFIG] smtphost      = [%s]", cfg.SMTP.Host)
	log.Printf("[CONFIG] smtpport      = [%d]", cfg.SMTP.Port)
	log.Printf("[CONFIG] smtpto        = [%s]", cfg.SMTP.To)
	log.Printf("[CONFIG] stubsmto      = [%t]", cfg.SMTP.DevMode)

	if cfg.DBHost == "" || cfg.DBUser == "" ||
		cfg.DBPass == "" || cfg.DBName == "" {
		flag.Usage()
		log.Printf("FATAL: Missing DB configuration")
		os.Exit(1)
	}
}
