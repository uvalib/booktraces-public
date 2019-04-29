package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

	log.Printf("Authentication successful for %s", computingID)
	json, _ := json.Marshal(user)
	if devMode {
		c.SetCookie("bt_admin_user", string(json), 3600, "/", "", false, false)
	} else {
		c.SetCookie("bt_admin_user", string(json), 3600, "/", "", true, false)
	}
	c.Redirect(http.StatusFound, "/admin")
}

// AuthMiddleware sits in front of all admin API calls and makes sure the shibboleth info is present
func (svc *ServiceContext) AuthMiddleware(c *gin.Context) {
	log.Printf("AuthMiddleware is checking for access cookie")
	cookieStr, err := c.Cookie("bt_admin_user")
	if err != nil {
		log.Printf("ERROR: unable to retrieve access cookie")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	log.Printf("AuthMiddleware found cookie %s; verifying", cookieStr)
	cookieUser := User{}
	err = json.Unmarshal([]byte(cookieStr), &cookieUser)
	if err != nil {
		log.Printf("ERROR: Unable to parse cookie: %s", err.Error())
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	user := User{}
	err = user.FindByEmail(svc.DB, cookieUser.Email)
	if err != nil {
		log.Printf("No user record found for %s. Not authorized.", cookieUser.Email)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if cookieUser.Token != user.Token {
		log.Printf("User access token mismatch")
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	log.Printf("User %s is authorized for %s", user.Email, c.Request.RequestURI)
	c.Next()
}
