package main

import (
	"fmt"
	"time"

	"context"

	"github.com/jianhan/ms-bmp-crawler/crawlers"
	"github.com/jianhan/ms-bmp-crawler/outputs"
	pcategories "github.com/jianhan/ms-bmp-products/proto/categories"
	pproducts "github.com/jianhan/ms-bmp-products/proto/products"
	psuppliers "github.com/jianhan/ms-bmp-products/proto/suppliers"
	cfgreader "github.com/jianhan/pkg/configs"
	"github.com/micro/go-micro"
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

	// crawl umart
	umart := crawlers.NewUmart(true)
	err = umart.Scrape()
	if err != nil {
		panic(err)
	}

	// output to service
	if err := serviceOutput.Output(context.Background(), umart); err != nil {
		panic(err)
	}

	// crawl megabuy
	megabuyau := crawlers.NewMegabuyau(true)
	megabuyau.Scrape()

	// output to service
	if err := serviceOutput.Output(context.Background(), megabuyau); err != nil {
		panic(err)
	}

}

func init() {
	viper.SetDefault("ENVIRONMENT", "development")
}
