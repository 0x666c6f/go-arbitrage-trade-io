package main

import (
	"flag"
	"fmt"
	"github.com/florianpautot/go-arbitrage-trade-io/arbitrage"
	"github.com/florianpautot/go-arbitrage-trade-io/http"
	"github.com/florianpautot/go-arbitrage-trade-io/tradeio"
	"github.com/florianpautot/go-arbitrage-trade-io/utils"
	"github.com/golang/glog"
	"os"
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
	// This is wa
	flag.Parse()
}

func main() {

	config,err := utils.LoadConfig("config.yaml")
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	http.Config = config
	tradeio.Config = config
	arbitrage.Config = config

	infos,err := tradeio.Info()
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	arbitrage.Infos = utils.FormatInfos(infos.Symbols)

	startDate := time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),time.Now().Hour(),time.Now().Minute()+1,config.StartSecond,0,time.Local)
	glog.Info("Start defined at ",startDate)
	sleep := startDate.Sub(time.Now())
	glog.Info("Starting arbitrages in ",sleep)
	time.Sleep(startDate.Sub(time.Now()))
	arbitrage.Start()
}