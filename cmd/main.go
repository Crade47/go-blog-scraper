package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	// "github.com/gocolly/colly"
	// "github.com/Crade47/medium-blog-scraper/internal/scraper"
	"github.com/Crade47/medium-blog-scraper/internal/scraper"
	"github.com/Crade47/medium-blog-scraper/models"
	"github.com/Crade47/medium-blog-scraper/utils"
	"github.com/rs/xid"
)

func main() {

	ctx := context.Background()

	//------------------------------INIT THE DATABASE------------------------------
	initErr := utils.Initialize(ctx, "config/firebase-config.json")
	if initErr != nil {
		fmt.Println("Error in initializing firebase")
		log.Fatal(initErr)
	}

	//------------------------------HANDLER FUNCTIONS------------------------------
	h1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("frontend/pages/index.html"))
		tmpl.Execute(w, nil)
	}

	createHandler := func(w http.ResponseWriter, r *http.Request) {
		id := xid.New().String()
		url := r.PostFormValue("medium-url")
		err := utils.CheckMediumURL(url)
		firestoreClient := utils.GetFirestoreClient()
		if err != nil {
			errObj := &models.APIError{
				Error: "Invalid URL.",
			}
			fmt.Println("URL error", err)
			// utils.WriteJSON(w, http.StatusBadRequest, errObj)
			tmpl := template.Must(template.ParseFiles("frontend/pages/index.html"))
			tmpl.Execute(w, errObj)
		} else {
			_, err := scraper.MediumScraper(ctx, id, url)
			if err != nil {
				errObj := &models.APIError{
					Error: "Server Error, Try again later.",
				}
				log.Fatal(err)
				tmpl := template.Must(template.ParseFiles("frontend/pages/index.html"))
				tmpl.Execute(w, errObj)
				return
			}
			blogData, err := firestoreClient.GetDocument(ctx, id)
			if err != nil {
				errObj := &models.APIError{
					Error: "Server Error, Try again later.",
				}
				log.Fatal(err)
				tmpl := template.Must(template.ParseFiles("frontend/pages/index.html"))
				tmpl.Execute(w, errObj)
				return
			}
			tmpl := template.Must(template.ParseFiles("frontend/pages/scraped_result.html"))
			tmpl.Execute(w, blogData)

		}
	}

	fs := http.FileServer(http.Dir("frontend"))

	http.HandleFunc("/", h1)
	http.HandleFunc("/create", createHandler)
	http.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	//------------------------------STARTING THE SERVER------------------------------
	fmt.Println("Listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))

	//SCRAPING THE BLOG
	// Blog, err := scraper.MediumScraper(ctx, id, "https://medium.com/blogging-guide/understanding-html-elements-of-medium-post-82d7e4b54826")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// //UPDATE THE DATABASE
	// firestoreClient := utils.GetFirestoreClient()
	// firestoreClient.AddDocument(ctx, id, Blog)

}
