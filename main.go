package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

const (
	filename = "product.csv"
)

type Product struct {
	URL         string
	Title       string
	Price       string
	Description string
	Image       string
}

func main() {
	var mu sync.Mutex
	var pagesMutex sync.Mutex
	var wg sync.WaitGroup

	var products []Product
	var visitedPages = make(map[string]bool)

	c := colly.NewCollector(
		colly.AllowedDomains("www.web-scraping.dev"),
		colly.Async(true),
	)

	c.OnHTML(".product", func(e *colly.HTMLElement) {
		product := Product{
			URL:         e.ChildAttr("h3 a", "href"),
			Image:       e.ChildAttr("img", "src"),
			Title:       e.ChildText("h3"),
			Price:       e.ChildText("div.price"),
			Description: e.ChildText("div.short-description"),
		}

		mu.Lock()
		products = append(products, product)
		mu.Unlock()
	})

	c.OnHTML(".paging a[href]", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
		if strings.Contains(nextPage, "products?page=") {
			absURL := e.Request.AbsoluteURL(nextPage)

			pagesMutex.Lock()
			if !visitedPages[absURL] {
				visitedPages[absURL] = true
				pagesMutex.Unlock()

				wg.Add(1)
				go func(url string) {
					defer wg.Done()
					err := c.Visit(url)
					if err != nil {
						log.Printf("Error while visiting %v", err)
					}
				}(absURL)
			} else {
				pagesMutex.Unlock()
				log.Printf("Page already visited: %s", absURL)
			}
		}
	})

	visitedPages["https://www.web-scraping.dev/products"] = true

	err := c.Visit("https://www.web-scraping.dev/products")
	if err != nil {
		log.Fatalf("Error visiting base url: %v", err)
	}

	c.Wait()
	wg.Wait()
	saveToCSV(products)
}

func saveToCSV(products []Product) {

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln("Error while creating CSV file:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{
		"Url",
		"Image",
		"Title",
		"Price",
		"Description",
	}
	if err := writer.Write(headers); err != nil {
		log.Fatalln("Error writing headers:", err)
	}

	for _, product := range products {
		record := []string{
			product.URL,
			product.Image,
			product.Title,
			product.Price,
			product.Description,
		}
		if err := writer.Write(record); err != nil {
			log.Printf("Error writing product: %v", err)
		}
	}

}
