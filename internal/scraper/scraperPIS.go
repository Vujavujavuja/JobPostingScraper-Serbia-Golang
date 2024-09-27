package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"job-aggregator/internal/models"
	"log"
	"strconv"
	"strings"
)

// ScrapeJobsPIS scrapes job listings from Poslovi Infostud
func ScrapeJobsPIS(baseURL string) []models.Job {
	var allJobs []models.Job
	var idAutoincrement int = 1
	jobsPerPage := 30

	// Step 1: Extract total number of jobs
	totalJobs := extractTotalJobs(baseURL)
	if totalJobs == 0 {
		log.Println("Failed to extract total number of jobs, exiting.")
		return allJobs
	}

	// Calculate total number of pages needed
	totalPages := (totalJobs / jobsPerPage) + 1 // Add one to cover partial pages
	fmt.Printf("Total jobs: %d, Total pages to scrape: %d\n", totalJobs, totalPages)

	// Step 2: Scrape each page
	for page := 1; page <= totalPages; page++ {
		// Build the URL for the current page
		url := fmt.Sprintf("%s?page=%d", baseURL, page)
		fmt.Println("Scraping URL:", url)

		// Initialize a new collector
		c := colly.NewCollector()

		var jobs []models.Job

		// On every job element
		c.OnHTML(".job", func(e *colly.HTMLElement) {
			job := models.Job{
				ID:        idAutoincrement,
				Title:     e.ChildText("h2.uk-margin-remove-bottom.uk-text-break"), // Job title
				Company:   e.ChildText("p.job-company"),                            // Company name
				Location:  e.ChildText("p.job-location"),                           // Job location
				Seniority: setSeniorityPIS(e.ChildText("span[data-tag-id='13']"), e.ChildText("h2.uk-margin-remove-bottom.uk-text-break")),
				URL:       "https://poslovi.infostud.com" + e.ChildAttr("a[href]", "href"), // Job URL
				Site:      "Poslovi Infostud",
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

// extractTotalJobs scrapes the total number of jobs from the job listings page
func extractTotalJobs(baseURL string) int {
	var totalJobs int

	// Initialize a new collector
	c := colly.NewCollector()

	// Extract the total number of jobs from the span element
	c.OnHTML("span.uk-text-muted.uk-margin-remove", func(e *colly.HTMLElement) {
		jobCountText := e.Text // e.g. "(3890 rezultata)"

		// Extract the digits from the text (assume the format is "(XXXX rezultata)")
		jobCountStr := strings.TrimPrefix(jobCountText, "(")
		jobCountStr = strings.TrimSuffix(jobCountStr, " rezultata)")

		// Convert job count string to an integer
		totalJobs, _ = strconv.Atoi(jobCountStr)
	})

	// Visit the base URL to extract total jobs
	err := c.Visit(baseURL)
	if err != nil {
		log.Printf("Failed to scrape base URL %s: %v", baseURL, err)
	}

	return totalJobs
}

func setSeniorityPIS(seniority string, jobTitle string) string {
	if seniority == "" {
		if strings.Contains(jobTitle, "internal") || strings.Contains(jobTitle, "Internal") || strings.Contains(jobTitle, "International") || strings.Contains(jobTitle, "international") {
			if strings.Contains(jobTitle, "junior") || strings.Contains(jobTitle, "trainee") || strings.Contains(jobTitle, "Junior") {
				return "Junior"
			} else if strings.Contains(jobTitle, "senior") || strings.Contains(jobTitle, "lead") || strings.Contains(jobTitle, "Senior") || strings.Contains(jobTitle, "Manager") || strings.Contains(jobTitle, "manager") || strings.Contains(jobTitle, "Head") || strings.Contains(jobTitle, "head") || strings.Contains(jobTitle, "Director") || strings.Contains(jobTitle, "director") || strings.Contains(jobTitle, "Menadzer") || strings.Contains(jobTitle, "Menadzer") {
				return "Senior"
			} else if strings.Contains(jobTitle, "mid") || strings.Contains(jobTitle, "middle") || strings.Contains(jobTitle, "intermediate") || strings.Contains(jobTitle, "medior") {
				return "Mid"
			} else {
				return "Mid"
			}
		}
		if strings.Contains(jobTitle, "junior") || strings.Contains(jobTitle, "trainee") || strings.Contains(jobTitle, "Junior") {
			return "Junior"
		} else if strings.Contains(jobTitle, "senior") || strings.Contains(jobTitle, "lead") || strings.Contains(jobTitle, "Senior") || strings.Contains(jobTitle, "Manager") || strings.Contains(jobTitle, "manager") || strings.Contains(jobTitle, "Head") || strings.Contains(jobTitle, "head") || strings.Contains(jobTitle, "Director") || strings.Contains(jobTitle, "director") || strings.Contains(jobTitle, "Menadzer") || strings.Contains(jobTitle, "Menadzer") {
			return "Senior"
		} else if strings.Contains(jobTitle, "mid") || strings.Contains(jobTitle, "middle") || strings.Contains(jobTitle, "intermediate") || strings.Contains(jobTitle, "medior") {
			return "Mid"
		} else if strings.Contains(jobTitle, "intern") || strings.Contains(jobTitle, "Intern") || strings.Contains(jobTitle, "praksa") || strings.Contains(jobTitle, "Praksa") || strings.Contains(jobTitle, "praktikant") || strings.Contains(jobTitle, "Praktikant") || strings.Contains(jobTitle, "internship") || strings.Contains(jobTitle, "Internship") {
			return "Intern"
		} else {
			return "Mid"
		}
	} else {
		if strings.Contains(jobTitle, "internal") || strings.Contains(jobTitle, "Internal") || strings.Contains(jobTitle, "International") || strings.Contains(jobTitle, "international") {
			if strings.Contains(jobTitle, "junior") || strings.Contains(jobTitle, "trainee") || strings.Contains(jobTitle, "Junior") {
				return "Junior"
			} else if strings.Contains(jobTitle, "senior") || strings.Contains(jobTitle, "lead") || strings.Contains(jobTitle, "Senior") || strings.Contains(jobTitle, "Manager") || strings.Contains(jobTitle, "manager") || strings.Contains(jobTitle, "Head") || strings.Contains(jobTitle, "head") || strings.Contains(jobTitle, "Director") || strings.Contains(jobTitle, "director") || strings.Contains(jobTitle, "Menadzer") || strings.Contains(jobTitle, "Menadzer") {
				return "Senior"
			} else if strings.Contains(jobTitle, "mid") || strings.Contains(jobTitle, "middle") || strings.Contains(jobTitle, "intermediate") || strings.Contains(jobTitle, "medior") {
				return "Mid"
			} else {
				return "Mid"
			}
		}
		if strings.Contains(seniority, "junior") || strings.Contains(seniority, "trainee") || strings.Contains(seniority, "Junior") {
			return "Junior"
		} else if strings.Contains(seniority, "senior") || strings.Contains(seniority, "lead") || strings.Contains(seniority, "Senior") {
			return "Senior"
		} else if strings.Contains(seniority, "mid") || strings.Contains(seniority, "middle") || strings.Contains(seniority, "intermediate") || strings.Contains(seniority, "medior") {
			return "Mid"
		} else if strings.Contains(seniority, "intern") || strings.Contains(seniority, "Intern") {

		} else {
			return "Mid"
		}
	}
	return "Mid"
}
