package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocolly/colly/v2"
)

type Hyperlink struct {
	Href string
	Text string
}

// http://localhost:8080/search?url=http://paulgraham.com/articles.html

func ping(w http.ResponseWriter, r *http.Request) {
	log.Println("Ping")
	w.Write([]byte("ping"))
}

func main() {

	addr := ":8080"

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/search", getLinks)
	http.HandleFunc("/readingList", getContents)
	http.HandleFunc("/ping", ping)
	log.Println("listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//Verify the param "URL" exists
	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing URL argument")
		return
	}
	log.Println("visiting", URL)

	//Create a new collector which will be in charge of collect the data from HTML
	c := colly.NewCollector()

	//Slices to store the data
	var response []Hyperlink

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		text := e.Text

		hl := Hyperlink{
			Href: link,
			Text: text,
		}

		if link != "" && text != "" {
			response = append(response, (hl))
		}
	})

	c.Visit(URL)

	object, err := json.Marshal(response)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}
	// Add some header and write the body for our endpoint
	w.Header().Add("Content-Type", "application/json")
	w.Write(object)
}

func getContents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// TODO: add get request to /search?url=PG_BLOG and loop through all the links
	// TODO: serve it as a mobile/kindle friendly webpage

	//Verify the param "URL" exists
	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing URL argument")
		return
	}
	log.Println("visiting", URL)

	//Create a new collector which will be in charge of collect the data from HTML
	c := colly.NewCollector()

	//Slices to store the data
	var response []string

	c.OnHTML("font", func(e *colly.HTMLElement) {
		text := e.Text

		if text != "" {
			response = append(response, text)
		}
	})

	c.Visit(URL)

	object, err := json.Marshal(response)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}
	// Add some header and write the body for our endpoint
	w.Header().Add("Content-Type", "application/json")
	w.Write(object)
}
