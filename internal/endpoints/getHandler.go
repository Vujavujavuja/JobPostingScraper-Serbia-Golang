package endpoints

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"job-aggregator/internal/models"
	"log"
	"net/http"
)

func GetJobsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jobs []models.Job

		// Query to select all jobs from the database
		rows, err := db.Query("SELECT id, title, company, location, seniority, url, site FROM jobs")
		if err != nil {
			log.Println("Failed to query jobs:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
			return
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Println("Failed to close rows:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close rows"})
				return
			}
		}(rows)

		// Iterate through the rows and scan each one into the Job struct
		for rows.Next() {
			var job models.Job
			if err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Location, &job.Seniority, &job.URL, &job.Site); err != nil {
				log.Println("Failed to scan job:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan job"})
				return
			}
			jobs = append(jobs, job)
		}

		// Respond with the list of jobs as JSON
		c.JSON(http.StatusOK, jobs)
	}
}
