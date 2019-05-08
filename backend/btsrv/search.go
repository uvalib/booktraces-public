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
	log.Printf("Search for [%s]", c.Query("q"))
	searchStr := fmt.Sprintf("%%%s%%", c.Query("q"))
	q := svc.DB.NewQuery(`select s.id,s.title, group_concat(distinct t.name  separator ", ") as tags,
		CONCAT('/uploads/', f.filename) as url,
		description, DATE_FORMAT(submitted_at,'%M %Y') as submitted
		from submissions s
		inner join submission_files f on f.submission_id = s.id
		left outer  join submission_tags st on st.submission_id = s.id
		left outer  join tags t on t.id = st.tag_id
		where public=1
			and (t.name like {:q} or s.title like {:q} or s.author like {:q}
	 		or s.publication_info like {:q} or s.library like {:q} or s.call_number like {:q}
	 		or s.description like {:q} or s.description like {:q} or s.submitted_at like {:q})
		group by s.id
		order by submitted_at`)
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
