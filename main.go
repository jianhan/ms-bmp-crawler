package main

import (
	"fmt"
	"time"

	"context"

	"sync"

	"github.com/jianhan/ms-bmp-crawler/crawlers"
	"github.com/jianhan/ms-bmp-crawler/outputs"
	pcategories "github.com/jianhan/ms-bmp-products/proto/categories"
	pproducts "github.com/jianhan/ms-bmp-products/proto/products"
	psuppliers "github.com/jianhan/ms-bmp-products/proto/suppliers"
	cfgreader "github.com/jianhan/pkg/configs"
	"github.com/micro/go-micro"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	serviceConfigs, err := cfgreader.NewReader(viper.GetString("ENVIRONMENT")).Read()
	if err != nil {
		panic(fmt.Sprintf("error while reading configurations: %s", err.Error()))
	}

	// initialize new service
	srv := micro.NewService(
		micro.Name(serviceConfigs.Name),
		micro.RegisterTTL(time.Duration(serviceConfigs.RegisterTTL)*time.Second),
		micro.RegisterInterval(time.Duration(serviceConfigs.RegisterInterval)*10),
		micro.Version(serviceConfigs.Version),
		micro.Metadata(serviceConfigs.Metadata),
	)

	// init service
	srv.Init()

	// initialize output
	serviceOutput := outputs.NewService(
		pcategories.NewCategoriesServiceClient("", nil),
		pproducts.NewProductsServiceClient("", nil),
		psuppliers.NewSuppliersServiceClient("", nil),
	)

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
	}()
	logrus.Info("Finish all scraping")
	wg.Wait()

	// output to service
	if err := serviceOutput.Output(context.Background(), umart); err != nil {
		panic(err)
	}

	// output to service
	if err := serviceOutput.Output(context.Background(), megabuyau); err != nil {
		panic(err)
	}

}

func init() {
	viper.SetDefault("ENVIRONMENT", "development")
}
