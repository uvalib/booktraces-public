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
	if svc.DevAuthUser != "" {
		computingID = svc.DevAuthUser
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
	log.Printf(string(json))
	c.SetCookie("bt_admin_user", string(json), 3600, "/", "", false, false)
	c.Redirect(http.StatusFound, "/admin")
}

// AuthMiddleware sits in front of all admin API calls and makes sure the shibboleth info is present
func (svc *ServiceContext) AuthMiddleware(c *gin.Context) {
	log.Printf("AuthMiddleware is checking for access cookie")
	cookieStr, err := c.Cookie("bt_admin_user")
	if err != nil {
		log.Printf("ERROR: unable to retrieve access cookie")
		c.Redirect(http.StatusForbidden, "/forbidden")
		return
	}
	log.Printf("AuthMiddleware found cookie %s; verifying", cookieStr)
	cookieUser := User{}
	err = json.Unmarshal([]byte(cookieStr), &cookieUser)
	if err != nil {
		log.Printf("ERROR: Unable to parse cookie: %s", err.Error())
		c.Redirect(http.StatusForbidden, "/forbidden")
		return
	}

	user := User{}
	err = user.FindByEmail(svc.DB, cookieUser.Email)
	if err != nil {
		log.Printf("No user record found for %s. Not authorized.", cookieUser.Email)
		c.Redirect(http.StatusForbidden, "/forbidden")
		return
	}

	if cookieUser.Token != user.Token {
		log.Printf("User access token mismatch")
		c.Redirect(http.StatusForbidden, "/forbidden")
		return
	}

	log.Printf("User %s is authorized for %s", user.Email, c.Request.RequestURI)
	c.Header("cache-control", "private, max-age=0, no-cache")

	c.Next()
}
