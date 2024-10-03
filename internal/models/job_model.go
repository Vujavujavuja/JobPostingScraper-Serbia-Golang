package models

// Job represents a job listing
// TODO: Figure out how to generate keywords for each job (myb gpt1.0) hugging face api is good too

type Job struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Company   string `json:"company"`
	Location  string `json:"location"`
	Seniority string `json:"seniority"`
	URL       string `json:"url"`
	Site      string `json:"site"`
}
