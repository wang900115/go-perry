package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

type Item struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	PubDate     time.Time `xml:"pubDate"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	PubDate     time.Time `xml:"pubDate"`
	Items       []Item    `xml:"item"`
}

func generateRSSFeed() ([]byte, error) {
	items := []Item{
		{
			Title:       "Article 1",
			Link:        "http://example.com/article1",
			Description: "Description of Article 1",
			PubDate:     time.Now(),
		},
		{
			Title:       "Article 2",
			Link:        "http://example.com/article2",
			Description: "Description of Article 2",
			PubDate:     time.Now().Add(-1 * time.Hour),
		},
		{
			Title:       "Article 3",
			Link:        "http://example.com/article3",
			Description: "Description of Article 3",
			PubDate:     time.Now(),
		},
	}

	feed := Channel{
		Title:       "Sample RSS Feed",
		Link:        "http://example.com",
		Description: "This is a sample RSS feed",
		PubDate:     time.Now(),
		Items:       items,
	}

	xmlData, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return nil, err
	}

	rssFeed := []byte(xml.Header + string(xmlData))
	return rssFeed, nil

}

func main() {

	rssFeed, err := generateRSSFeed()
	if err != nil {
		fmt.Println("Error generating RSS feed:", err)
		return
	}

	file, err := os.Create("feed.xml")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(rssFeed)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("RSS feed generated successfully and saved to feed.xml")
}
