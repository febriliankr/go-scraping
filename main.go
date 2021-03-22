package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocolly/colly/v2"
)

// http://localhost:8080/search?url=http://paulgraham.com/articles.html

func ping(w http.ResponseWriter, r *http.Request) {
	log.Println("Ping")
	w.Write([]byte("ping"))
}

func main() {

	addr := ":8080"

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/search", getData)
	http.HandleFunc("/ping", ping)
	log.Println("listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

type hyperlink struct {
	href string
	text string
}

func getData(w http.ResponseWriter, r *http.Request) {
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
	var response []hyperlink

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		text := e.Text

		hl := hyperlink{
			href: link,
			text: text,
		}

		if link != "" && text != "" {
			response = append(response, (hl))
		}
	})

	c.Visit(URL)

	log.Println(response)

	object, err := json.Marshal(response)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}
	// Add some header and write the body for our endpoint
	w.Header().Add("Content-Type", "application/json")
	w.Write(object)
}
