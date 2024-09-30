package models

// Job represents a job listing
// TODO: Figure out how to generate keywords for each job (myb gpt1.0)
type Job struct {
	ID        int
	Title     string
	Company   string
	Location  string
	Seniority string
	URL       string
	Site      string
}
