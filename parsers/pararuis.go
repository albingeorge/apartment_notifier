package parsers

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type pararius struct {
	url string
}

func (p *pararius) GetApartments() ([]string, error) {
	var apartments []string
	resp, err := http.Get((*p).url)

	if err != nil {
		log.Println("HTTP fetch failed: ", (*p).url)
		// IF http fetch failed, we should return immediately, else the defer statement will fail
		// Alternative is to return err as well and handle gracefully in the calling function
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Invalid resp code: ", resp.StatusCode, "; url: ", (*p).url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Println("error parsing webpage ", (*p).url)
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

	return apartments, nil
}
