package main

import (
	"context"
	"fmt"
	"time"

	"github.com/albingeorge/pararius_notifier/parsers"
	"github.com/martinlindhe/notify"
)

func main() {
	go notifyNewApartments("https://www.pararius.com/apartments/amsterdam/0-1600", "Amsterdam")

	go notifyNewApartments("https://www.pararius.com/apartments/utrecht/0-1500", "Utrecht")

	go notifyNewApartments("https://www.pararius.com/apartments/amersfoort/0-1500", "Amersfoot")

	go notifyNewApartments("https://www.pararius.com/apartments/haarlem/0-1500/50m2", "Haarlem")

	// go notifyNewApartments("https://www.funda.nl/en/huur/amsterdam,utrecht,diemen,amersfoort/1000-1750/1-dag/sorteer-datum-af/", "Amsterdam Funda")
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
	parser, err := parsers.GetParser(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	newApartments, err := parser.GetApartments()
	if err != nil {
		fmt.Println(err)
		return
	}

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
