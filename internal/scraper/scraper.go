package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"job-aggregator/internal/models"
	"log"
	"reflect"
)

// ScrapeJobs scrapes job listings from a website with pagination
func ScrapeJobsPIS(baseURL string) []models.Job {
	var allJobs []models.Job
	page := 1
	var firstPageJobs []models.Job

	for {
		// Build the URL for the current page
		url := fmt.Sprintf("%s&page=%d", baseURL, page)
		fmt.Println("Scraping URL:", url)

		// Initialize a new collector
		c := colly.NewCollector()

		var jobs []models.Job

		// On every job element
		c.OnHTML(".job", func(e *colly.HTMLElement) {
			job := models.Job{
				Title:     e.ChildText("h2.uk-margin-remove-bottom.uk-text-break"), // Job title
				Company:   e.ChildText("p.job-company"),                            // Company name
				Location:  e.ChildText("p.job-location"),                           // Job location
				Type:      e.ChildText(".job-type"),                                // Job type
				Seniority: e.ChildText("span[data-tag-id='13']"),                   // Seniority
				URL:       e.ChildAttr("a[href]", "href"),                          // Job URL
			}

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

		// Check if the jobs scraped from the current page are identical to the first page
		if page == 1 {
			// Save the first page jobs for future comparison
			firstPageJobs = jobs
		} else if reflect.DeepEqual(firstPageJobs, jobs) {
			// If the jobs from this page are the same as the first page, stop pagination
			fmt.Println("Encountered duplicate jobs from page 1, stopping pagination.")
			break
		}

		// Append the jobs from this page to the total list
		allJobs = append(allJobs, jobs...)

		// Increment the page number and scrape the next page
		page++
	}

	return allJobs
}
