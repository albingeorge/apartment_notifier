package parsers

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type funda struct {
	url string
}

func (f *funda) GetApartments() ([]string, error) {
	var apartments []string

	doc, err := getDocumentFromUrl((*f).url)

	if err != nil {
		return nil, err
	}

	doc.Find(".search-result-content .search-result__header-title").
		Each(func(i int, s *goquery.Selection) {
			fmt.Println("Getting funda page: " + s.Text())
			apartments = append(apartments, s.Text())
		})

	return apartments, nil
}
