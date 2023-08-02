package scraper

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/Crade47/medium-blog-scraper/models"
	"github.com/Crade47/medium-blog-scraper/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func UploadHTML(ctx context.Context, destination string) (string, error) {
	storageClientObj := utils.GetFirebaseStorage()
	formattedDestination := destination + ".html"
	htmlUrl, err := storageClientObj.UploadFile(ctx, "output.html", formattedDestination)

	if err != nil {
		return "", nil
	}

	if err := os.Remove("output.html"); err != nil {
		return "", nil
	}

	return htmlUrl, err
}

func MediumScraper(ctx context.Context, id string, url string) (*models.Blog, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("medium.com"),
	)
	// CREATING A FILE
	file, err := os.OpenFile("output.html", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	//------------------------------CALL BACK FUNCTIONS------------------------------

	var newBlog *models.Blog

	//Assigning ID
	newBlog.ID = id

	//------------------------------CALL BACK FUNCTIONS------------------------------

	c.OnHTML("h1.pw-post-title", func(e *colly.HTMLElement) {

		if newBlog == nil {
			newBlog = &models.Blog{
				Title: e.Text,
			}
		} else {
			newBlog.Title = e.Text
		}

	})

	c.OnHTML("div.bl p.be a.af", func(e *colly.HTMLElement) {
		link := "medium.com" + e.Attr("href")

		if newBlog == nil {
			newBlog = &models.Blog{
				AuthorName: e.Text,
				AuthorURL:  link,
			}
		} else {

			newBlog.AuthorName = e.Text
			newBlog.AuthorURL = link
		}

	})

	// META TAGS
	c.OnHTML("meta", func(e *colly.HTMLElement) {

		name := e.Attr("name")
		property := e.Attr("property")

		// FILTERING THE USEFUL SEO META TAGS
		if name == "description" || name == "keywords" || property == "og:title" || property == "og:description" ||
			name == "author" || property == "article:published_time" || property == "article:modified_time" ||
			property == "article:tag" || property == "article:section" || property == "og:url" ||
			name == "twitter:creator" || property == "article:author" {

			html, err := goquery.OuterHtml(e.DOM)
			if err != nil {
				log.Fatalln(err)
			}

			_, writeErr := file.Write([]byte(html))
			if err != nil {
				fmt.Println("Error writing to the file")
				log.Fatalln(writeErr)
			}
		}

	})

	c.OnHTML(" h1.aej, h2, p.pw-post-body-paragraph, ul, img.bg", func(e *colly.HTMLElement) {
		classRegex := regexp.MustCompile(`class\s*=\s*["']([^"']*)["']`)
		switch e.Name {
		case "p":
			html, err := e.DOM.Html()
			if err != nil {
				log.Fatalln(err)
			}

			data := "<p class='paragraph'>" + html + "</p>\n"

			_, writeErr := file.Write([]byte(data))
			if err != nil {
				fmt.Println("Error writing to the file")
				log.Fatalln(writeErr)
			}
		case "h1":
			html, err := e.DOM.Html()
			if err != nil {
				log.Fatalln(err)
			}

			data := "<h1 class='heading'>" + html + "</h1>\n"
			_, writeErr := file.Write([]byte(data))
			if err != nil {
				fmt.Println("Error writing to the file")
				log.Fatalln(writeErr)
			}
		case "h2":
			html, err := e.DOM.Html()
			if err != nil {
				log.Fatalln(err)
			}

			data := "<h2 class='subHeading'>" + html + "</h2>\n"
			_, writeErr := file.Write([]byte(data))
			if err != nil {
				fmt.Println("Error writing to the file")
				log.Fatalln(writeErr)
			}
		case "ul":
			html, err := e.DOM.Html()
			if err != nil {
				log.Fatalln(err)
			}

			html = classRegex.ReplaceAllString(html, "class='list-item'")
			data := "<ul class='unordered-list'>" + html + "</ul>"
			_, writeErr := file.Write([]byte(data))
			if writeErr != nil {
				fmt.Println("Error writing to the file")
				log.Fatalln(writeErr)
			}
		case "blockquote":
			html, err := e.DOM.Html()
			if err != nil {
				log.Fatalln(err)
			}
			data := "<blockquote class='blockquote'>" + html + "</blockquote>\n"
			_, writeErr := file.Write([]byte(data))
			if err != nil {
				fmt.Println("Error writing to the file")
				log.Fatalln(writeErr)
			}
		case "img":
			e.DOM.SetAttr("class", "blog-picture")
			outerHtml, err := goquery.OuterHtml(e.DOM)
			if err != nil {
				log.Fatalln(err)
			}
			_, writeErr := file.Write([]byte(outerHtml))
			if err != nil {
				fmt.Println("Error writing to the file")
				log.Fatalln(writeErr)
			}
		}

	})

	visitErr := c.Visit(url)
	if err != nil {
		return nil, visitErr
	}

	//------------------------------CALLING UPLOAD FUNCTION------------------------------
	htmlUrl, err := UploadHTML(ctx, id)
	if err != nil {
		return nil, err
	}

	newBlog.HtmlFileURL = htmlUrl

	return newBlog, nil
}
