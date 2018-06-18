package crawlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	pcategories "github.com/jianhan/ms-bmp-products/proto/categories"
	pproducts "github.com/jianhan/ms-bmp-products/proto/products"
)

type base struct {
	name        string
	categoryURL string
	categories  []*pcategories.Category
	products    []*pproducts.Product
	currency    string
	homepageURL string
	testMode    bool
}

func (b *base) Name() string {
	return b.name
}

func (b *base) Categories() []*pcategories.Category {
	return b.categories
}

func (b *base) Products() []*pproducts.Product {
	return b.products
}

func (b *base) addCategory(c *pcategories.Category) {
	b.categories = append(b.categories, c)
}

func (b *base) addProduct(p *pproducts.Product) {
	b.products = append(b.products, p)
}

func (b *base) htmlDoc(url string) (*goquery.Document, func() error, error) {
	// get html page
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		return nil, nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, nil, err
	}

	return doc, res.Body.Close, nil
}

func (b *base) priceStrToFloat(priceStr string) (price float64, err error) {
	priceRaw := strings.Replace(priceStr, " ", "", -1)
	priceRaw = strings.Replace(priceRaw, ",", "", -1)
	priceRaw = strings.Replace(priceRaw, "$", "", -1)
	price, err = strconv.ParseFloat(priceRaw, 64)
	if err != nil {
		return
	}

	return
}

func (b *base) getLinkFullURL(url string) string {
	if strings.HasPrefix(url, b.homepageURL) {
		return url
	}
	url = strings.Replace(url, " ", "", -1)
	url = strings.Trim(url, "/")

	return strings.Trim(b.homepageURL, "/") + "/" + url
}

func (b *base) HomepageURL() string {
	return b.homepageURL
}

func (b *base) Currency() string {
	return b.currency
}

func (b *base) Validate() error {
	if strings.Trim(b.Name(), " ") == "" {
		return errors.New("empty name")
	}

	if len(b.Categories()) == 0 {
		return errors.New("empty categories")
	}

	if len(b.Products()) == 0 {
		return errors.New("empty products")
	}

	return nil
}
