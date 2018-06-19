package main

import (
	"fmt"
	"time"

	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/jianhan/ms-bmp-crawler/crawlers"
	pcategories "github.com/jianhan/ms-bmp-products/proto/categories"
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

	umart := crawlers.NewUmart(true)
	umart.Scrape()
	categoriesClient := pcategories.NewCategoriesServiceClient("", nil)

	categoriesRsp, _ := categoriesClient.Categories(context.Background(), &pcategories.CategoriesReq{})
	umartCategories := umart.Categories()
	for _, v := range categoriesRsp.Categories {
		for k := range umartCategories {
			if umartCategories[k].Url == v.Url {
				umartCategories[k].ID = v.ID
			}
		}
	}
	rsp, err := categoriesClient.UpsertCategories(context.Background(), &pcategories.UpsertCategoriesReq{Categories: umartCategories})
	if err != nil {
		panic(err)
	}
	spew.Dump(rsp)

	//megabuyau := crawlers.NewMegabuyau(true)
	//megabuyau.Scrape()
	//spew.Dump(megabuyau.Products())

}

func init() {
	viper.SetDefault("ENVIRONMENT", "development")
}
