package main

import (
	"fmt"
	"job-aggregator/internal/scraper"
)

func main() {
	// URL of the job listings page you want to scrape
	url := "https://poslovi.infostud.com/oglasi-za-posao?category%5B0%5D=5" // Replace with actual job board URL

	// Call the ScrapeJobs function to fetch job listings
	jobs := scraper.ScrapeJobsPIS(url)

	// Print the results to check the scraper
	fmt.Println("Scraped Jobs:")
	for _, job := range jobs {
		fmt.Printf("ID: %d\n", job.ID)
		fmt.Printf("Title: %s\n", job.Title)
		fmt.Printf("Company: %s\n", job.Company)
		fmt.Printf("Location: %s\n", job.Location)
		fmt.Printf("Type: %s\n", job.Type)
		fmt.Printf("Seniority: %s\n", job.Seniority)
		fmt.Printf("URL: %s\n", job.URL)
		fmt.Println("-------------")
	}
}
