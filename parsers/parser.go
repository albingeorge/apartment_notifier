package parsers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func getDocumentFromUrl(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)

	if err != nil {
		log.Println("HTTP fetch failed: ", url)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("Invalid resp code: " + strconv.Itoa(resp.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, errors.New("error parsing webpage " + url)
	}

	return doc, nil
}
