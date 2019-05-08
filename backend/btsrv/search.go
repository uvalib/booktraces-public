package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

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
	if c.Query("t") != "" {
		getTaggedSubmissions(svc.DB, c, c.Query("t"))
		return
	}
	c.String(http.StatusBadRequest, "Invalid search type")
}

func getArchives(db *dbx.DB, c *gin.Context, tgtYearMonth string) {
	log.Printf("Get archives from [%s]", tgtYearMonth)
}

func getTaggedSubmissions(db *dbx.DB, c *gin.Context, tgtTag string) {
	log.Printf("Get submission tagged [%s]", tgtTag)
}

func doQuery(db *dbx.DB, c *gin.Context, query string) {
	log.Printf("Search for string [%s]", query)
	searchStr := fmt.Sprintf("%%%s%%", query)
	searchQ := fmt.Sprintf(`%s 
		and (t.name like {:q} or s.title like {:q} or s.author like {:q}
			or s.publication_info like {:q} or s.library like {:q} or s.call_number like {:q}
			or s.description like {:q} or s.description like {:q} or s.submitted_at like {:q})
		group by s.id
		order by submitted_at`, getBaseQuery())
	q := db.NewQuery(searchQ)
	q.Bind(dbx.Params{"q": searchStr})

	type Hit struct {
		ID          int     `json:"id" db:"id"`
		Title       string  `json:"title" db:"title"`
		Tags        *string `json:"tags" db:"tags"`
		URL         string  `json:"url" db:"url"`
		Description string  `json:"description" db:"description"`
		SubmittedOn string  `json:"submittedOn" db:"submitted"`
	}
	var hits []Hit
	err := q.All(&hits)
	if err != nil {
		log.Printf("ERROR: Unable to get search results: %s", err.Error())
		c.String(http.StatusNotFound, err.Error())
	}
	for idx := range hits {
		hit := &hits[idx]
		thumb := GetThumbFilename(hit.URL)
		log.Print(thumb)
		hit.URL = thumb
	}
	c.JSON(http.StatusOK, hits)
}

func getBaseQuery() string {
	return `select s.id,s.title, 
		group_concat(distinct t.name  separator ", ") as tags,
		CONCAT('/uploads/', f.filename) as url,
		description, DATE_FORMAT(submitted_at,'%M %Y') as submitted
		from submissions s
			inner join submission_files f on f.submission_id = s.id
			left outer  join submission_tags st on st.submission_id = s.id
			left outer  join tags t on t.id = st.tag_id
		where public=1 `
}
