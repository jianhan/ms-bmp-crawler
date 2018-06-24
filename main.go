package main

import (
	"context"
	"time"

	"context"

	"github.com/spf13/viper"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Minute)
	defer cancel()

	//serviceConfigs, err := cfgreader.NewReader(viper.GetString("ENVIRONMENT")).Read()
	//if err != nil {
	//	panic(fmt.Sprintf("error while reading configurations: %s", err.Error()))
	//}
	//
	//// initialize new service
	//srv := micro.NewService(
	//	micro.Name(serviceConfigs.Name),
	//	micro.RegisterTTL(time.Duration(serviceConfigs.RegisterTTL)*time.Second),
	//	micro.RegisterInterval(time.Duration(serviceConfigs.RegisterInterval)*10),
	//	micro.Version(serviceConfigs.Version),
	//	micro.Metadata(serviceConfigs.Metadata),
	//)
	//
	//// init service
	//srv.Init()

}

func init() {
	viper.SetDefault("ENVIRONMENT", "development")
}
