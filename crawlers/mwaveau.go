package crawlers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
)

type mwaveau struct {
	base
}

func NewMwaveau(testMode bool) Scraper {
	// TODO: this one has async loading menu, can not do it for now
	b := base{
		homepageURL: "https://www.mwave.com.au/",
		name:        "Mwave Australia",
		categoryURL: "https://www.mwave.com.au/",
		currency:    "AUD",
		testMode:    testMode,
	}

	return &mwaveau{b}
}

func (m *mwaveau) Scrape() error {
	// start scraping
	if err := m.fetchCategories(); err != nil {
		return err
	}
	if err := m.fetchProducts(); err != nil {
		return err
	}

	return nil
}

func (m *mwaveau) fetchCategories() error {
	// link text define all link with specific text that we want
	linkText := []string{
		"Components & Parts",
	}
	spew.Dump(linkText)
	doc, fn, err := m.htmlDoc(m.categoryURL)
	if err != nil {
		return err
	}
	defer fn()

	// get all links with class categoryLink
	doc.Find("li.hasChild").Each(func(i int, s *goquery.Selection) {
		spew.Dump("testtest")
		if linkNode := s.Find("a b").First().Get(0); linkNode != nil {
			// found it, get all categories within hover menu
			s.Find("a").Each(func(j int, linkS *goquery.Selection) {
				spew.Dump(linkS.Attr("href"))
			})
		}
	})

	return nil
}

func (m *mwaveau) fetchProducts() error {
	return nil
}
