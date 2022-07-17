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
	go notifyNewApartments("https://www.pararius.com/apartments/amsterdam/0-1500", "Amsterdam")

	go notifyNewApartments("https://www.pararius.com/apartments/utrecht/0-1500", "Utrecht")

	go notifyNewApartments("https://www.pararius.com/apartments/amersfoort/0-1500", "Amersfoot")
	select {}
}

func notifyNewApartments(url string, area string) {
	ctx := context.Background()

	compareApartments(ctx, url, area)

	for range time.Tick(time.Minute * 2) {
		compareApartments(ctx, url, area)
	}
}

func compareApartments(ctx context.Context, url string, area string) {
	// Get stored top apartment
	latestApartment := getInstance().redis.HGet(ctx, "latest_apartments", url).Val()
	fmt.Println("\nLATEST apartment: ", latestApartment)

	// Get current apartment listing
	newApartments := getNewApartments(url)

	// Check if currently stored apartment is at the top of listing
	if len(newApartments) > 0 && newApartments[0] != latestApartment {
		// Print diff if not
		fmt.Println("New apartments found!")
		printLatestApartments(latestApartment, newApartments, area)

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
		// IF http fetch failed, we should return immediately, else the defer statement will fail
		// Alternative is to return err as well and handle gracefully in the calling function
		return []string{}
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Invalid resp code: ", resp.StatusCode, "; url: ", url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Println("error parsing webpage ", url)
	}

	doc.Find(".search-list__item--listing").
		Each(func(i int, s *goquery.Selection) {
			// Ignore featured items
			if s.Find(".listing-label--featured").Size() == 0 {
				s.Find("a.listing-search-item__link--title").Each(func(i int, s *goquery.Selection) {
					apartments = append(apartments, s.Text())
				})
			}
		})

	return apartments
}

func printLatestApartments(oldLatestApartment string, newApartments []string, area string) {
	notifyMessage := ""
	fmt.Println("New apartments:")
	for _, apartment := range newApartments {
		if apartment == oldLatestApartment {
			break
		}
		notifyMessage = notifyMessage + apartment + "\n"
		fmt.Println(apartment)
	}
	notify.Alert("Apartment notifier", "New apartment found in "+area, notifyMessage, "path/to/icon.png")
}
