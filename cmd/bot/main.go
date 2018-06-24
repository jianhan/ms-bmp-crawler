package main

import (
	"context"
	"sync"
	"time"

	"github.com/jianhan/ms-bmp-crawler/crawlers"
	"github.com/jianhan/ms-bmp-crawler/outputs"
	pcategories "github.com/jianhan/ms-bmp-products/proto/categories"
	pproducts "github.com/jianhan/ms-bmp-products/proto/products"
	psuppliers "github.com/jianhan/ms-bmp-products/proto/suppliers"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Minute)
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(2)
	// crawl umart
	umart := crawlers.NewUmart(false)
	go func() {
		if err := umart.Scrape(); err != nil {
			wg.Done()
			panic(err)
		}
		wg.Done()
	}()

	megabuyau := crawlers.NewMegabuyau(false)
	go func() {
		// crawl megabuy
		if err := megabuyau.Scrape(); err != nil {
			wg.Done()
			panic(err)
		}
		wg.Done()
	}()

	wg.Wait()
	logrus.Info("Finish all scraping")

	// initialize output
	serviceOutput := outputs.NewService(
		pcategories.NewCategoriesServiceClient("", nil),
		pproducts.NewProductsServiceClient("", nil),
		psuppliers.NewSuppliersServiceClient("", nil),
	)

	// output to service
	if err := serviceOutput.Output(ctx, umart); err != nil {
		panic(err)
	}

	// output to service
	if err := serviceOutput.Output(ctx, megabuyau); err != nil {
		panic(err)
	}
}
