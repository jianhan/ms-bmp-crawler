package crawlers

import (
	"github.com/PuerkitoBio/goquery"
	pcategories "github.com/jianhan/ms-bmp-products/proto/categories"
	pproducts "github.com/jianhan/ms-bmp-products/proto/products"
	psuppliers "github.com/jianhan/ms-bmp-products/proto/suppliers"
	"github.com/sirupsen/logrus"
)

type umart struct {
	base
}

func NewUmart(testMode bool) Crawler {
	b := base{
		name:        "Umart",
		categoryURL: "https://www.umart.com.au/all-categories.html",
		homepageURL: "https://www.umart.com.au",
		currency:    "AUD",
		testMode:    testMode,
	}

	return &umart{b}
}

func (u *umart) Scrape() error {
	if err := u.fetchCategories(); err != nil {
		return err
	}
	if err := u.fetchProducts(); err != nil {
		return err
	}

	return nil
}

func (u *umart) fetchCategories() error {
	doc, fn, err := u.htmlDoc(u.categoryURL)
	if err != nil {
		return err
	}
	defer fn()

	// get all links with class categoryLink
	doc.Find("div.ovhide.productsIn.productText > a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok && href != "" {
			u.addCategory(&pcategories.Category{Name: s.Text(), Url: u.getLinkFullURL(href)})
		}
	})

	return nil
}

func (u *umart) fetchProducts() error {
	for _, c := range u.categories {
		if err := u.fetchProductsByURL(c.Url, c.Url); err != nil {
			return err
		}

		// test mode checking
		if u.testMode && len(u.products) > 0 {
			break
		}
	}

	return nil
}

func (u *umart) fetchProductsByURL(url, categoryURL string) error {
	doc, fn, err := u.htmlDoc(url)
	if err != nil {
		return err
	}
	defer fn()

	// find products
	doc.Find("li.goods_info").Each(func(i int, s *goquery.Selection) {
		p := &pproducts.Product{CategoryUrl: categoryURL, Currency: u.currency}

		// find image
		s.First().Find("div.goods_img > a > img").Each(func(imgI int, imgS *goquery.Selection) {
			src, ok := imgS.Attr("src")
			if ok {
				p.ImageUrl = src
			}
		})

		// find product name
		s.First().Find("div.content_holder1 > div.goods_name > a").Each(func(nameI int, nameS *goquery.Selection) {
			// product url
			href, ok := nameS.Attr("href")
			if ok {
				p.Url = href
			}

			// product name
			p.Name = nameS.Text()

		})

		// find product price
		s.First().Find("span.goods_price").Each(func(priceI int, priceS *goquery.Selection) {
			if p.Price, err = u.priceStrToFloat(priceS.Text()); err != nil {
				logrus.Warn(err)
			}
		})
		if p.Price > 0 {
			u.addProduct(p)
		}
	})

	// find next page url
	var nextPageURL string
	doc.Find("ul.page li a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == ">" {
			href, ok := s.Attr("href")
			if ok {
				nextPageURL = u.getLinkFullURL(href)
			}
		}
	})
	if nextPageURL != "" {
		u.fetchProductsByURL(nextPageURL, categoryURL)
	}

	return nil
}

func (u *umart) Supplier() *psuppliers.Supplier {
	return &psuppliers.Supplier{
		Name:        u.name,
		LogoUrl:     "https://assets.umart.com.au/themes/umart2018/images/logo_lg.png",
		HomePageUrl: u.homepageURL,
		Currency:    "AUD",
	}
}
