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
func ScrapeJobsHW(baseURL string) []models.Job {
	var allJobs []models.Job
	var idAutoincrement int = 1
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

		c := colly.NewCollector()

		var jobs []models.Job

	}

	return nil
}

func extractTotalJobsHW(baseURL string) int {
	var totalJobs int

	// Initialize a new collector
	c := colly.NewCollector()

	// Add debugging for the element and the selector
	c.OnHTML("h2 span", func(e *colly.HTMLElement) {
		// Log the content of each matched element
		fmt.Println("Matched element text:", e.Text)

		// Look for the span that contains the number of jobs
		jobCountText := e.Text
		if strings.Contains(jobCountText, "oglas") {
			// Clean up the job count text by removing unnecessary characters
			jobCountStr := strings.TrimPrefix(jobCountText, "(")      // Remove the opening parenthesis
			jobCountStr = strings.TrimSuffix(jobCountStr, " oglasa)") // Handle "oglasa)"
			jobCountStr = strings.TrimSuffix(jobCountStr, " oglas)")  // Handle "oglas)"

			// Try to convert the cleaned string to an integer
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
