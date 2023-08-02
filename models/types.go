package models

import "time"

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Blog struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	AuthorURL   string    `json:"author_url"`
	AuthorName  string    `json:"author_name"`
	PublishedOn time.Time `json:"published_on"`
	HtmlFileURL string    `json:"html_file_url"`
}
