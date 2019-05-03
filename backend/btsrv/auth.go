package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// Authenticate will authenticate a user based on Shibboleth headers, then redirect to other pages for content
// NOTE: this is called directly from the front-end as a transient page with  window.location.href = "/authenticate"
// to force it through NetBadge authentication. This route does nothing more than ensure users have been thru
// authentication and create a user if one does not exist. End result is a redirect.
func (svc *ServiceContext) Authenticate(c *gin.Context) {
	log.Printf("Checking authentication headers...")
	computingID := c.GetHeader("remote_user")
	devMode := false
	if svc.DevAuthUser != "" {
		computingID = svc.DevAuthUser
		devMode = true
		log.Printf("Using dev auth user ID: %s", computingID)
	}
	if computingID == "" {
		log.Printf("ERROR: Expected auth header not present in request. Not authorized.")
		c.Redirect(http.StatusFound, "/forbidden")
		return
	}

	email := fmt.Sprintf("%s@virginia.edu", computingID)
	user := User{}
	err := user.FindByEmail(svc.DB, email)
	if err != nil {
		log.Printf("No user record found for authorized computing ID. Not authorized. %s", err.Error())
		c.Redirect(http.StatusFound, "/forbidden")
		return
	}

	log.Printf("Adding access token to user")
	user.Token = xid.New().String()
	svc.DB.Model(&user).Update("Token")

	tgtURL := c.Query("url")
	log.Printf("Authentication successful for %s at URL %s", computingID, tgtURL)
	json, _ := json.Marshal(user)

	// Set user account into in an open cookie, and place the access token
	// in a secure, http-only cookie that the browser can't touch. it will just be
	// passed along on all admin api requests
	c.SetCookie("bt_admin_user", string(json), 3600, "/", "", false, false)
	adminSess := fmt.Sprintf("%s|%s", user.Token, user.Email)
	if devMode {
		c.SetCookie("bt_admin_session", adminSess, 0, "/", "", false, true)
	} else {
		c.SetCookie("bt_admin_session", adminSess, 0, "/", "", true, true)
	}
	c.Redirect(http.StatusFound, tgtURL)
}

// AuthMiddleware sits in front of all admin API calls and makes the auth token generated
// by the shibboleth-fronted authenticate handler is present and valid
func (svc *ServiceContext) AuthMiddleware(c *gin.Context) {
	log.Printf("AuthMiddleware is checking for access cookie")
	cookieStr, err := c.Cookie("bt_admin_session")
	if err != nil {
		log.Printf("ERROR: unable to retrieve access cookie")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	log.Printf("AuthMiddleware found cookie %s; verifying", cookieStr)
	user := User{}
	err = user.FindByToken(svc.DB, strings.Split(cookieStr, "|")[0])
	if err != nil {
		log.Printf("No user record found for token. Not authorized.")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if user.Email != strings.Split(cookieStr, "|")[1] {
		log.Printf("Email / token mismatch. Not authorized.")
		c.AbortWithStatus(http.StatusForbidden)
	}

	log.Printf("User %s is authorized for %s", user.Email, c.Request.RequestURI)
	c.Next()
}
