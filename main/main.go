package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"job-aggregator/internal/models"
	"job-aggregator/internal/scraper"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type OpenAiRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type OpenAiResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func main() {
	// Connect to SQLite database
	db, err := sql.Open("sqlite3", "./jobs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	// Create the jobs table if not exists
	createTable := `
	CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		company TEXT,
		location TEXT,
		seniority TEXT,
		url TEXT UNIQUE,
	    Site TEXT
	);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	// Delete all from jobs table before inserting new jobs
	deleteSQL := `
	DELETE FROM jobs;
	`

	_, err = db.Exec(deleteSQL)
	if err != nil {
		log.Fatal(err)
	}

	urlPIS := "https://poslovi.infostud.com/oglasi-za-posao/"
	urlHW := "https://www.helloworld.rs/oglasi-za-posao/"
	//urlPRS := "https://www.poslovi.rs/jobs"

	// Call the ScrapeJobs function to fetch job listings
	jobs := scraper.ScrapeJobsPIS(urlPIS)
	idAfterPIS := len(jobs) + 1
	jobs = append(jobs, scraper.ScrapeJobsHW(urlHW, idAfterPIS)...)
	//idAfterHW := len(jobs) + 1
	//jobs = append(jobs, scraper.ScrapeJobsPRS(urlPRS, idAfterHW)...)
	//jobs := scraper.ScrapeJobsHW(url)

	// Insert scraped jobs into the database
	for _, job := range jobs {
		err := insertJob(db, job)
		if err != nil {
			log.Printf("Failed to insert job '%s': %v\n", job.Title, err)
		}
	}

	fmt.Println("Jobs inserted into the database successfully.")

	queryJobs(db)

	// API functionality
	r := gin.Default()

	// GET route for all jobs
	r.GET("/jobs", func(c *gin.Context) {
		var jobs []models.Job

		// Query to select all jobs from the database
		rows, err := db.Query("SELECT id, title, company, location, seniority, url, site FROM jobs")
		if err != nil {
			log.Println("Failed to query jobs:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
			return
		}
		defer rows.Close()

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

		// Response
		c.JSON(http.StatusOK, jobs)
	})

	// run
	err = r.Run(":8080")
	if err != nil {
		return
	}
}

// insertJob inserts a job entry into the SQLite database
func insertJob(db *sql.DB, job models.Job) error {
	insertSQL := `
	INSERT OR IGNORE INTO jobs (title, company, location, seniority, url, site) 
	VALUES (?, ?, ?, ?, ?, ?);
	`
	_, err := db.Exec(insertSQL, job.Title, job.Company, job.Location, job.Seniority, job.URL, job.Site)
	return err
}

// queryJobs retrieves and displays jobs from the SQLite database (optional for verification)
func queryJobs(db *sql.DB) {
	rows, err := db.Query("SELECT id, title, company, location, seniority, url, site FROM jobs")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Jobs in the database:")
	for rows.Next() {
		var id int
		var title, company, location, jobType, seniority, url string
		err = rows.Scan(&id, &title, &company, &location, &jobType, &seniority, &url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d\nTitle: %s\nCompany: %s\nLocation: %s\nType: %s\nSeniority: %s\nURL: %s\n\n",
			id, title, company, location, jobType, seniority, url)
	}
}
