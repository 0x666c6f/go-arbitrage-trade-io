package main

import (
	"context"
	"flag"
	"github.com/adshao/go-binance"
	"github.com/florianpautot/go-arbitrage/arbitrage"
	"github.com/florianpautot/go-arbitrage/global"
	"github.com/florianpautot/go-arbitrage/utils"
	"github.com/golang/glog"
	"os"
	"strconv"
	"time"
)

func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "INFO")
	if len(os.Getenv("LogLevel")) > 0{
		flag.Set("v", os.Getenv("LogLevel"))
	}

	// This is wa
	flag.Parse()
}

func main() {

	loadedConfig,err := utils.LoadConfig("config.yaml")
	if err != nil {
		glog.V(1).Info(err.Error())
		return
	}

	if len(os.Getenv("StartSecond")) > 0{
		startSecond,err := strconv.Atoi(os.Getenv("StartSecond"))
		if err != nil {
			glog.V(1).Info(err.Error())
			return
		}
		loadedConfig.StartSecond = startSecond
	}

	if len(os.Getenv("EndSecond")) > 0 {
		endSecond,err := strconv.Atoi(os.Getenv("EndSecond"))
		if err != nil {
			glog.V(1).Info(err.Error())
			return
		}

		loadedConfig.EndSecond = endSecond
	}

	global.GlobalConfig = loadedConfig


	utils.UpdateCachedBalances()

	var (
		apiKey = global.GlobalConfig.APIKey
		secretKey = global.GlobalConfig.APISecret
	)
	global.Binance = binance.NewClient(apiKey, secretKey)

	infos,err := global.Binance.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		glog.V(1).Info(err.Error())
		return
	}
	arbitrage.Infos = utils.FormatInfos(infos.Symbols)


	startDate := time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),time.Now().Hour(),time.Now().Minute()+1, global.GlobalConfig.StartSecond,0,time.Local)
	glog.V(1).Info("Start defined at ",startDate)
	sleep := startDate.Sub(time.Now())
	glog.V(1).Info("Starting arbitrages in ",sleep)
	time.Sleep(startDate.Sub(time.Now()))
	arbitrage.Start()
}

