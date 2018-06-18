package crawlers

import (
	pcategories "github.com/jianhan/ms-bmp-products/proto/categories"
	pproducts "github.com/jianhan/ms-bmp-products/proto/products"
)

type Scraper interface {
	Name() string
	Scrape() error
	Categories() []*pcategories.Category
	Products() []*pproducts.Product
	Validate() error
	HomepageURL() string
	Currency() string
}
