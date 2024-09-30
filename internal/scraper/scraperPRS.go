package scraper

import (
	"bytes"
	"fmt"
	"github.com/gocolly/colly"
	"io"
	"job-aggregator/internal/models"
	"log"
	"net/http"
	"strconv"
)

// ScrapeJobsPRS scrapes job listings from Poslovi.rs using AJAX requests DOESNT WORK CURRENTLY
// TODO: Fix the scraper to work with Poslovi.rs
func ScrapeJobsPRS(baseURL string, firstId int) []models.Job {
	var allJobs []models.Job
	var idAutoincrement int = firstId
	page := 1
	jobsPerPage := 30

	// Scrape each page of jobs using AJAX
	for {
		// Adjust the `filter_start` to simulate pagination
		filterStart := strconv.Itoa((page - 1) * jobsPerPage)

		// Initialize a new collector
		c := colly.NewCollector()

		var jobs []models.Job

		// On every job element (adjust the selectors based on the actual HTML structure)
		c.OnHTML("div.single_job_listing", func(e *colly.HTMLElement) {
			job := models.Job{
				ID:        idAutoincrement,
				Title:     e.ChildText("h3 a"),                   // Job title
				Company:   e.ChildText("div.company"),            // Company name
				Location:  e.ChildText("div.location"),           // Job location
				Seniority: "Unknown",                             // Seniority (if available)
				URL:       baseURL + e.ChildAttr("h3 a", "href"), // Job URL
				Site:      "Poslovi.rs",
			}
			idAutoincrement++
			jobs = append(jobs, job)
		})

		// Handle errors
		c.OnError(func(r *colly.Response, err error) {
			log.Printf("Failed to scrape %s with error: %s", r.Request.URL, err)
		})

		// Create the AJAX request with the POST data (mimicking the browser's AJAX request)
		ajaxURL := baseURL + "jobs_ajax/90"
		data := fmt.Sprintf("is_list=1&search_keywords=&filter_start=%s&filter_stop=0&filter_sort=1", filterStart)
		headers := map[string]string{
			"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
			"X-Requested-With": "XMLHttpRequest",
			"User-Agent":       "Mozilla/5.0",
		}

		// Perform the AJAX request
		err := postRequest(ajaxURL, data, headers)
		if err != nil {
			log.Printf("Error making POST request: %v", err)
			break
		}

		// Visit the AJAX response URL
		err = c.Visit(ajaxURL)
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

		// Increment the page number for the next iteration
		page++
	}

	return allJobs
}

// postRequest performs a POST request with custom data and headers
func postRequest(url string, data string, headers map[string]string) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}

	// Set custom headers for the POST request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check if the response was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to make POST request: %s", resp.Status)
	}

	return nil
}
