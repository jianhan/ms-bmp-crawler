package crawlers

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	pcategories "github.com/jianhan/ms-bmp-products/proto/categories"
	pproducts "github.com/jianhan/ms-bmp-products/proto/products"
	"github.com/sirupsen/logrus"
)

type megabuyau struct {
	base
}

func NewMegabuyau(testMode bool) Crawler {
	b := base{
		homepageURL: "https://www.megabuy.com.au",
		name:        "Megabuy Australia",
		categoryURL: "https://www.megabuy.com.au/computer-components-c1160.html",
		currency:    "AUD",
		testMode:    testMode,
	}

	return &megabuyau{b}
}

func (m *megabuyau) Scrape() error {
	// start scraping
	if err := m.fetchCategories(); err != nil {
		return err
	}
	if err := m.fetchProducts(); err != nil {
		return err
	}

	return nil
}

func (m *megabuyau) fetchCategories() error {
	doc, fn, err := m.htmlDoc(m.categoryURL)
	if err != nil {
		return err
	}
	defer fn()

	// get all links with class categoryLink
	doc.Find("a.categoryLink").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		href, ok := s.Attr("href")
		if ok && href != "" {
			m.addCategory(&pcategories.Category{Name: s.Text(), Url: href})
		}
	})

	return nil
}

func (m *megabuyau) fetchProducts() error {
	for _, c := range m.categories {
		if err := m.fetchProductsByURL(c.Url, c.Url); err != nil {
			return err
		}

		// test mode checking
		if m.testMode && len(m.products) > 0 {
			break
		}
	}

	return nil
}

func (m *megabuyau) fetchProductsByURL(url, categoryURL string) error {
	doc, fn, err := m.htmlDoc(url)
	if err != nil {
		return err
	}
	defer fn()

	// find products
	doc.Find("div.productListing div.productListingRow, div.productListing div.productListingRowAlt").Each(func(i int, s *goquery.Selection) {
		p := &pproducts.Product{CategoryUrl: categoryURL, Currency: m.currency}

		// find image
		s.Find("div.image > a > img").First().Each(func(ii int, is *goquery.Selection) {
			src, ok := is.Attr("src")
			if ok {
				p.ImageUrl = src
			}
		})

		// find name
		s.Find("div.nameDescription > a").First().Each(func(ni int, ns *goquery.Selection) {
			p.Name = ns.Text()
		})

		// find price
		s.Find("div.price > span").First().Each(func(ni int, ns *goquery.Selection) {
			if p.Price, err = m.priceStrToFloat(ns.Text()); err != nil {
				logrus.Warn(err)
			}
		})
		if p.Price > 0 {
			// append product into products
			m.addProduct(p)
		}
	})

	// find next page url
	doc.Find("div.pagination").First().Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(ai int, as *goquery.Selection) {
			title, ok := as.Attr("title")
			if ok {
				if strings.ToLower(strings.Trim(title, " ")) == "next page" {
					nextPageHref, ok := as.Attr("href")
					if ok {
						m.fetchProductsByURL(nextPageHref, categoryURL)
					}
				}
			}
		})
	})

	return nil
}
