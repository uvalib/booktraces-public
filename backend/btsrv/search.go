package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

type searchHit struct {
	ID          int     `json:"id" db:"id"`
	Title       string  `json:"title" db:"title"`
	Tags        *string `json:"tags" db:"tags"`
	URL         string  `json:"url" db:"url"`
	Description string  `json:"description" db:"description"`
	SubmittedOn string  `json:"submittedOn" db:"submitted"`
	Institution string  `json:"institution" db:"institution"`
}

// Search executes a full search over the submissions table
func (svc *ServiceContext) Search(c *gin.Context) {
	if c.Query("q") != "" {
		doQuery(svc.DB, c, c.Query("q"))
		return
	}
	if c.Query("a") != "" {
		getArchives(svc.DB, c, c.Query("a"))
		return
	}
	if c.Query("i") != "" {
		instID, _ := strconv.Atoi(c.Query("i"))
		getByInstitution(svc.DB, c, instID)
		return
	}
	if c.Query("t") != "" {
		getTaggedSubmissions(svc.DB, c, c.Query("t"))
		return
	}
	c.String(http.StatusBadRequest, "Invalid search type")
}

func getByInstitution(db *dbx.DB, c *gin.Context, institutionID int) {
	log.Printf("INFO: get submissions by institution [%d]", institutionID)
	searchQ := fmt.Sprintf(`%s and institution_id={:i} group by s.id order by submitted_at`, getBaseQuery())
	q := db.NewQuery(searchQ)
	q.Bind(dbx.Params{"i": institutionID})
	getHits(c, q)
}

func getArchives(db *dbx.DB, c *gin.Context, tgtYearMonth string) {
	log.Printf("INFO: get archives from [%s]", tgtYearMonth)
	year := strings.Split(tgtYearMonth, "-")[0]
	month := strings.Split(tgtYearMonth, "-")[1]
	searchQ := fmt.Sprintf(`%s 
		and year(submitted_at) = {:y} and month(submitted_at) = {:m} 
		group by s.id
		order by submitted_at`, getBaseQuery())
	q := db.NewQuery(searchQ)
	q.Bind(dbx.Params{"y": year, "m": month})
	getHits(c, q)
}

func getTaggedSubmissions(db *dbx.DB, c *gin.Context, tgtTag string) {
	log.Printf("INFO: get submission tagged [%s]", tgtTag)
	searchQ := fmt.Sprintf(`%s group by s.id having find_in_set({:t}, tags) order by submitted_at`, getBaseQuery())
	q := db.NewQuery(searchQ)
	q.Bind(dbx.Params{"t": tgtTag})
	getHits(c, q)
}

func doQuery(db *dbx.DB, c *gin.Context, query string) {
	log.Printf("INFO: search for string [%s]", query)
	searchStr := fmt.Sprintf("%%%s%%", query)
	searchQ := fmt.Sprintf(`%s 
		and (t.name like {:q} or s.title like {:q} or s.author like {:q}
			or s.publication_info like {:q} or s.library like {:q} or s.call_number like {:q}
			or s.description like {:q} or s.description like {:q} or s.submitted_at like {:q}
			or s.submitter_name like {:q})
		group by s.id
		order by submitted_at`, getBaseQuery())
	q := db.NewQuery(searchQ)
	q.Bind(dbx.Params{"q": searchStr})
	getHits(c, q)
}

func getHits(c *gin.Context, q *dbx.Query) {
	var hits []searchHit
	err := q.All(&hits)
	if err != nil {
		log.Printf("ERROR: Unable to get search results: %s", err.Error())
		c.String(http.StatusNotFound, err.Error())
	}
	for idx := range hits {
		hit := &hits[idx]
		thumb := getThumbFilename(hit.URL)
		hit.URL = thumb
	}
	c.JSON(http.StatusOK, hits)
}

func getBaseQuery() string {
	return `select s.id,s.title,i.name as institution, 
		group_concat(distinct t.name  separator ", ") as tags,
		CONCAT('/uploads/', f.filename) as url,
		description, DATE_FORMAT(submitted_at,'%M %Y') as submitted
		from submissions s
			inner join submission_files f on f.submission_id = s.id
			inner join institutions i on i.id = s.institution_id
			left outer  join submission_tags st on st.submission_id = s.id
			left outer  join tags t on t.id = st.tag_id
		where public=1 `
}
