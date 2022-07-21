package parsers

import (
	"github.com/PuerkitoBio/goquery"
)

type pararius struct {
	url string
}

func (p *pararius) GetApartments() ([]string, error) {
	var apartments []string

	doc, err := getDocumentFromUrl((*p).url)

	if err != nil {
		return nil, err
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
