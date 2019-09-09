package main

import (
	"flag"
	"github.com/florianpautot/go-arbitrage-trade-io/arbitrage"
	"github.com/florianpautot/go-arbitrage-trade-io/http"
	"github.com/florianpautot/go-arbitrage-trade-io/tradeio"
	"github.com/florianpautot/go-arbitrage-trade-io/utils"
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

	config,err := utils.LoadConfig("config.yaml")
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
		config.StartSecond = startSecond

	}

	if len(os.Getenv("EndSecond")) > 0 {
		endSecond,err := strconv.Atoi(os.Getenv("EndSecond"))
		if err != nil {
			glog.V(1).Info(err.Error())
			return
		}

		config.EndSecond = endSecond
	}

	balances, err := tradeio.Account()
	if err != nil {
		glog.V(1).Info(err.Error())
	}

	if len(balances.Balances) > 0 {
		formattedBalances := utils.FormatBalance(balances.Balances)
		config.MaxBTC,err = strconv.ParseFloat(formattedBalances["btc"].Available,64)
		if err != nil {
			glog.V(1).Info(err.Error())
		}
		config.MaxUSDT,err = strconv.ParseFloat(formattedBalances["usdt"].Available,64)
		if err != nil {
			glog.V(1).Info(err.Error())
		}
		config.MaxETH,err = strconv.ParseFloat(formattedBalances["eth"].Available,64)
		if err != nil {
			glog.V(1).Info(err.Error())
		}
	}

	http.Config = config
	tradeio.Config = config
	arbitrage.Config = config

	infos,err := tradeio.Info()
	if err != nil {
		glog.V(1).Info(err.Error())
		return
	}
	arbitrage.Infos = utils.FormatInfos(infos.Symbols)


	startDate := time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),time.Now().Hour(),time.Now().Minute()+1,config.StartSecond,0,time.Local)
	glog.V(1).Info("Start defined at ",startDate)
	sleep := startDate.Sub(time.Now())
	glog.V(1).Info("Starting arbitrages in ",sleep)
	time.Sleep(startDate.Sub(time.Now()))
	arbitrage.Start()
}