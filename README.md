# Work in progress app
## Job Postings aggregator app for the Serbian job market
Contains a simple web scraper for the most popular job posting boards in Serbia. Aggregates them into a SQLite database.
Plans for future functionality include:
CV Matching;
Algorithmic job searches for matches;
Embedded HuggingFace GPT2 Wrapper or something similar as a companion.

# Job database model:
  ID        int
	Title     string
	Company   string
	Location  string
	Seniority string
	URL       string
	Site      string
