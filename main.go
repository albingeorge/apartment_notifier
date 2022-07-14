package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/martinlindhe/notify"
)

func main() {
	go notifyNewApartments("https://www.pararius.com/apartments/amsterdam/0-1500")

	go notifyNewApartments("https://www.pararius.com/apartments/utrecht/0-1500")
	select {}
}

func notifyNewApartments(url string) {
	ctx := context.Background()

	compareApartments(ctx, url)

	for range time.Tick(time.Minute * 2) {
		compareApartments(ctx, url)
	}
}

func compareApartments(ctx context.Context, url string) {
	// Get stored top apartment
	latestApartment := getInstance().redis.HGet(ctx, "latest_apartments", url).Val()
	fmt.Println("\nLATEST apartment: ", latestApartment)

	// Get current apartment listing
	newApartments := getNewApartments(url)

	// Check if currently stored apartment is at the top of listing
	if len(newApartments) > 0 && newApartments[0] != latestApartment {
		// Print diff if not
		fmt.Println("New apartments found!")
		printLatestApartments(latestApartment, newApartments)

		// Update the stored top apartment
		getInstance().redis.HSet(ctx, "latest_apartments", url, newApartments[0])
	} else {
		fmt.Println("No new apartments found!")
	}
}

func getNewApartments(url string) []string {
	var apartments []string
	resp, err := http.Get(url)

	if err != nil {
		log.Println("HTTP fetch failed: ", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Invalid resp code: ", resp.StatusCode, "; url: ", url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Println("error parsing webpage ", url)
	}

	doc.Find(".listing-search-item__title a").
		Each(func(i int, s *goquery.Selection) {
			// fmt.Println("Index: ", i)
			// fmt.Println("Name: ", s.Text())
			apartments = append(apartments, s.Text())
		})

	return apartments
}

func printLatestApartments(oldLatestApartment string, newApartments []string) {
	notifyMessage := "New apartments:\n"
	fmt.Println("New apartments:")
	for _, apartment := range newApartments {
		if apartment == oldLatestApartment {
			break
		}
		notifyMessage = notifyMessage + apartment + "\n"
		fmt.Println(apartment)
	}
	notify.Notify("Apartment notifier", "New apartment found", notifyMessage, "path/to/icon.png")
}
