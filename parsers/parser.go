package parsers

import (
	"errors"
	"strings"
)

type Parser interface {
	GetApartments() ([]string, error)
}

func GetParser(url string) (Parser, error) {
	if strings.HasPrefix(url, "https://www.pararius.com") {
		return &pararius{url: url}, nil
	} else if strings.HasPrefix(url, "https://www.funda.nl") {
		return &funda{url: url}, nil
	}

	return nil, errors.New("no parser found for the url: " + url)
}
