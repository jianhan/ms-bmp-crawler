package main

import (
	"fmt"
	"time"

	"github.com/jianhan/ms-bmp-crawler/crawlers"
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

	//umart := crawlers.NewUmart(true)
	//umart.Scrape()

	//megabuyau := crawlers.NewMegabuyau(true)
	//megabuyau.Scrape()
	//spew.Dump(megabuyau.Products())

	mwaveau := crawlers.NewMwaveau(true)
	mwaveau.Scrape()
}

func init() {
	viper.SetDefault("ENVIRONMENT", "development")
}
