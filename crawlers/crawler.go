package crawlers

import (
	pcategories "github.com/jianhan/ms-bmp-products/proto/categories"
	pproducts "github.com/jianhan/ms-bmp-products/proto/products"
	psuppliers "github.com/jianhan/ms-bmp-products/proto/suppliers"
)

type Crawler interface {
	Name() string
	Scrape() error
	Categories() []*pcategories.Category
	Products() []*pproducts.Product
	Validate() error
	HomepageURL() string
	Currency() string
	Supplier() *psuppliers.Supplier
}
