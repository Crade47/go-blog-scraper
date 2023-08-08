package models

import "time"

type Blog struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	AuthorURL    string    `json:"author_url"`
	AuthorName   string    `json:"author_name"`
	PublishedOn  time.Time `json:"published_on"`
	HtmlFileURL  string    `json:"html_file_url"`
	MainImageURL string    `json:"main_image_url"`
}

type APIError struct {
	Error string `json:"error"`
}
