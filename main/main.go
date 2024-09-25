package main

import (
	"database/sql"
	"fmt"
	"job-aggregator/internal/models"
	"job-aggregator/internal/scraper"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	// Connect to SQLite database (it will create the file if it doesn't exist)
	db, err := sql.Open("sqlite3", "./jobs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the jobs table if it doesn't exist
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

	// URL of the job listings page you want to scrape
	urlPIS := "https://poslovi.infostud.com/oglasi-za-posao?category%5B0%5D=5" // Poslovi InfoStud - IT jobs
	urlHW := "https://www.helloworld.rs/oglasi-za-posao/programiranje"

	// Call the ScrapeJobs function to fetch job listings
	jobs := scraper.ScrapeJobsPIS(urlPIS)
	idAfterPIS := len(jobs) + 1
	jobs = append(jobs, scraper.ScrapeJobsHW(urlHW, idAfterPIS)...)
	//jobs := scraper.ScrapeJobsHW(url)

	// Insert scraped jobs into the database
	for _, job := range jobs {
		err := insertJob(db, job)
		if err != nil {
			log.Printf("Failed to insert job '%s': %v\n", job.Title, err)
		}
	}

	fmt.Println("Jobs inserted into the database successfully.")

	// Optionally, retrieve and display the jobs from the database to verify
	queryJobs(db)
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
