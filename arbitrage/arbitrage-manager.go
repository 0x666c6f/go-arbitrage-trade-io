package arbitrage

import (
	"context"
	"github.com/adshao/go-binance"
	"github.com/florianpautot/go-arbitrage/global"
	"github.com/florianpautot/go-arbitrage/utils"
	"github.com/golang/glog"
	"strconv"
	"time"
)

var Infos map[string]binance.Symbol
var TotalMinuteWeight = 0
var TotalMinuteOrderWeight = 0

var valBTC float64
var valETH float64
var valEthBTC float64

func Start() {
	utils.UpdateCachedBalances()
	restartDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute()+1, global.GlobalConfig.StartSecond, 0, time.Local)
	glog.V(2).Info("Starting arbitrage")

	for TotalMinuteWeight < (global.GlobalConfig.APIMinuteLimit-23) && global.GlobalConfig.EndSecond > time.Now().Second() {
		launchArbitrages()
		if global.GlobalConfig.Timeout != "" {
			duration, err := time.ParseDuration(global.GlobalConfig.Timeout)
			if err != nil {
				glog.V(1).Info(err.Error())
			}
			time.Sleep(duration)
		}
	}

	utils.UpdateCachedBalances()

	TotalMinuteWeight = 0;
	TotalMinuteOrderWeight = 0;
	if time.Now().Before(restartDate) {
		sleepTime := restartDate.Sub(time.Now())
		glog.V(2).Info("Will sleep ", sleepTime.Seconds(), "to reset minute weight");
		time.Sleep(sleepTime)
		glog.V(2).Info("Waking up, sleep is over !");
	}
	Start()
}

func launchArbitrages() {
	tickers, err := global.Binance.NewListBookTickersService().Do(context.Background())
	TotalMinuteWeight += 20
	if err != nil {
		glog.V(1).Info(err.Error())
		glog.V(1).Info("Error while retrieving tickers, will sleep until next loop")
		wakeUp := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute()+2, global.GlobalConfig.StartSecond, 0, time.Local)
		glog.V(1).Info("Will sleep ", wakeUp.Sub(time.Now()))
		time.Sleep(wakeUp.Sub(time.Now()))
		glog.V(1).Info("Waking up, back to work !")
		return
	}

	formattedTickers, symbols := utils.FormatTickers(tickers)
	valBTC, err = strconv.ParseFloat(formattedTickers["btc_usdt"].AskPrice, 64)
	if err != nil {
		glog.V(1).Info(err.Error())
	}
	valETH, err = strconv.ParseFloat(formattedTickers["eth_usdt"].AskPrice, 64)
	if err != nil {
		glog.V(1).Info(err.Error())
	}
	valEthBTC, err = strconv.ParseFloat(formattedTickers["eth_btc"].AskPrice, 64)
	if err != nil {
		glog.V(1).Info(err.Error())
	}

	symbolsLen := len(symbols);
	for index := 0; index < symbolsLen; index++ {
		symbol := symbols[index]
		UsdtToBtcEthToUsdt(formattedTickers, Infos, symbol, "btc")
		BtcEthBtcArbitrage(formattedTickers, Infos, symbol)
		EthBtcToUsdtBtcToEthBtc(formattedTickers, Infos, symbol, "eth", "btc")
		EthBtcToUsdtBtcToEthBtc(formattedTickers, Infos, symbol, "btc", "usdt")
		EthBtcToUsdtBtcToEthBtc(formattedTickers, Infos, symbol, "eth", "usdt")
		UsdtToBtcEthToUsdt(formattedTickers, Infos, symbol, "eth")

	}
}
