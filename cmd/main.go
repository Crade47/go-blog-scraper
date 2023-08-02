package main

import (
	"context"
	"fmt"
	"log"

	// "github.com/gocolly/colly"
	"github.com/Crade47/medium-blog-scraper/internal/scraper"
	"github.com/Crade47/medium-blog-scraper/utils"
	"github.com/rs/xid"
)

func main() {
	id := xid.New().String()

	ctx := context.Background()
	//INIT THE DATABASE
	initErr := utils.Initialize(ctx, "config/firebase-config.json")
	if initErr != nil {
		fmt.Println("Error in initializing firebase")
		log.Fatal(initErr)
	}
	//SCRAPING THE BLOG
	Blog, err := scraper.MediumScraper(ctx, id, "https://medium.com/blogging-guide/understanding-html-elements-of-medium-post-82d7e4b54826")
	if err != nil {
		log.Fatalln(err)
	}

	//UPDATE THE DATABASE
	firestoreClient := utils.GetFirestoreClient()
	firestoreClient.AddDocument(ctx, id, Blog)

}
