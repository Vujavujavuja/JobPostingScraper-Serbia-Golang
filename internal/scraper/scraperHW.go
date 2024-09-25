package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"job-aggregator/internal/models"
	"log"
	"strconv"
	"strings"
)

// ScrapeJobsHW scrapes job listings from Helloworld.rs
func ScrapeJobsHW(baseURL string, firstId int) []models.Job {
	var allJobs []models.Job
	var idAutoincrement int = firstId
	pageModifier := 30
	jobsPerPage := 30

	// Step 1: Extract total number of jobs
	totalJobs := extractTotalJobsHW(baseURL)

	if totalJobs == 0 {
		fmt.Println("ERROR: Failed to extract total number of jobs, exiting.")
		return allJobs
	}

	totalPages := (totalJobs / jobsPerPage) + 1
	fmt.Printf("Total jobs: %d, Total pages to scrape: %d\n", totalJobs, totalPages)

	// Step 2: Scrape each page
	for page := 1; page <= totalPages; page++ {
		url := fmt.Sprintf("%s?page=%d", baseURL, (page-1)*pageModifier)
		fmt.Println("Scraping URL:", url)

		// Initialize a new collector
		c := colly.NewCollector()

		var jobs []models.Job

		// On every job element
		c.OnHTML("div.relative.z-1.flex.flex-col", func(e *colly.HTMLElement) {
			job := models.Job{
				ID:        idAutoincrement,
				Title:     e.ChildText("h3 a.__ga4_job_title"),                       // Job title
				Company:   e.ChildText("h4 a.__ga4_job_company"),                     // Company name
				Location:  cleanLoc(e.ChildText("p.text-sm.font-semibold")),          // Job location
				Seniority: e.ChildText("button.__ga4_job_seniority"),                 // Seniority
				URL:       "https://www.helloworld.rs" + e.ChildAttr("h3 a", "href"), // Job URL
				Site:      "Helloworld.rs",
			}
			idAutoincrement++
			jobs = append(jobs, job)
		})

		// Handle errors (if the page doesn't exist or the request fails)
		c.OnError(func(r *colly.Response, err error) {
			log.Printf("Failed to scrape %s with error: %s", r.Request.URL, err)
		})

		// Visit the current page
		err := c.Visit(url)
		if err != nil {
			log.Fatal(err)
		}

		// If no jobs are found on the current page, exit the loop
		if len(jobs) == 0 {
			fmt.Printf("No jobs found on page %d, stopping pagination.\n", page)
			break
		}

		// Append the jobs from this page to the total list
		allJobs = append(allJobs, jobs...)
	}

	return allJobs
}

// extractTotalJobsHW extracts the total number of jobs from the job listings page
func extractTotalJobsHW(baseURL string) int {
	var totalJobs int

	// Initialize a new collector
	c := colly.NewCollector()

	// Extract the total number of jobs from the h2 span element
	c.OnHTML("h2 span", func(e *colly.HTMLElement) {
		// Log the content of the matched element for debugging
		jobCountText := e.Text
		if strings.Contains(jobCountText, "oglas") {
			// Clean up the job count text by removing unnecessary characters
			jobCountStr := strings.TrimPrefix(jobCountText, "(")
			jobCountStr = strings.TrimSuffix(jobCountStr, " oglasa)")
			jobCountStr = strings.TrimSuffix(jobCountStr, " oglas)")

			// Convert job count string to an integer
			var err error
			totalJobs, err = strconv.Atoi(jobCountStr)
			if err != nil {
				log.Printf("Error converting job count to integer: %v", err)
			} else {
				fmt.Printf("Extracted job count: %d\n", totalJobs)
			}
		}
	})

	// Visit the base URL to extract total jobs
	err := c.Visit(baseURL)
	if err != nil {
		log.Printf("Failed to scrape base URL %s: %v", baseURL, err)
	}

	return totalJobs
}

func cleanLoc(location string) string {
	// check if the string contains a number and remove everything after it
	if strings.Contains(location, "1") {
		location = strings.Split(location, "1")[0]
	}
	if strings.Contains(location, "2") {
		location = strings.Split(location, "2")[0]
	}
	if strings.Contains(location, "3") {
		location = strings.Split(location, "3")[0]
	}
	if strings.Contains(location, "4") {
		location = strings.Split(location, "4")[0]
	}
	if strings.Contains(location, "5") {
		location = strings.Split(location, "5")[0]
	}
	if strings.Contains(location, "6") {
		location = strings.Split(location, "6")[0]
	}
	if strings.Contains(location, "7") {
		location = strings.Split(location, "7")[0]
	}
	if strings.Contains(location, "8") {
		location = strings.Split(location, "8")[0]
	}
	if strings.Contains(location, "9") {
		location = strings.Split(location, "9")[0]
	}
	if strings.Contains(location, "0") {
		location = strings.Split(location, "0")[0]
	}
	return location
}
