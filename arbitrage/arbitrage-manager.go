package arbitrage

import (
	"github.com/florianpautot/go-arbitrage-trade-io/model"
	"github.com/florianpautot/go-arbitrage-trade-io/model/responses"
	"github.com/florianpautot/go-arbitrage-trade-io/tradeio"
	"github.com/florianpautot/go-arbitrage-trade-io/utils"
	"github.com/golang/glog"
	"strconv"
	"time"
)

var Config model.Config
var Infos map[string]responses.Symbol
var TotalMinuteWeight = 0
var TotalMinuteOrderWeight = 0

var valBTC float64
var valETH float64
var valEthBTC float64

func Start(){
	restartDate := time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),time.Now().Hour(),time.Now().Minute()+1,Config.StartSecond,0,time.Local)
	glog.V(1).Info("Starting arbitrage")

	for TotalMinuteWeight < (Config.APIMinuteLimit - 23) && Config.EndSecond > time.Now().Second() {
		launchArbitrages()
		if Config.Timeout != ""{
			duration,err:= time.ParseDuration(Config.Timeout)
			if err != nil {
				glog.V(1).Info(err.Error())
			}
			time.Sleep(duration)
		}
	}

	balances, err := tradeio.Account()
	if err != nil {
		glog.V(1).Info(err.Error())
	}

	if len(balances.Balances) > 0 {
		formattedBalances := utils.FormatBalance(balances.Balances)
		Config.MaxBTC,err = strconv.ParseFloat(formattedBalances["btc"].Available,64)
		if err != nil {
			glog.V(1).Info(err.Error())
		}
		Config.MaxUSDT,err = strconv.ParseFloat(formattedBalances["usdt"].Available,64)
		if err != nil {
			glog.V(1).Info(err.Error())
		}
		Config.MaxETH,err = strconv.ParseFloat(formattedBalances["eth"].Available,64)
		if err != nil {
			glog.V(1).Info(err.Error())
		}
	}

	TotalMinuteWeight = 0;
	TotalMinuteOrderWeight = 0;
	if time.Now().Before(restartDate) {
		sleepTime := restartDate.Sub(time.Now())
		glog.V(1).Info("Will sleep", sleepTime.Seconds(), "to reset minute weight");
		time.Sleep(sleepTime)
		glog.V(1).Info("Waking up, sleep is over !");
	}
	Start()
}

func launchArbitrages(){
	glog.V(1).Info("Launching arbitrages")
	tickers,err := tradeio.Tickers()
	if err != nil{
		glog.V(1).Info(err.Error())
	}
	TotalMinuteWeight +=20
	if tickers.Code != 0 {
		glog.V(1).Info("Error while retrieving tickers, will sleep until next loop")
		wakeUp := time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),time.Now().Hour(),time.Now().Minute(),Config.StartSecond+1,0,time.Local)
		glog.V(1).Info("Will sleep",wakeUp)
		time.Sleep(wakeUp.Sub(time.Now()))
		glog.V(1).Info("Waking up, back to work !")
		return
	}

	formattedTickers,symbols := utils.FormatTickers(tickers.Tickers)
	valBTC,err = strconv.ParseFloat(formattedTickers["btc_usdt"].AskPrice,64)
	if err != nil{
		glog.V(1).Info(err.Error())
	}
	valETH,err = strconv.ParseFloat(formattedTickers["eth_usdt"].AskPrice,64)
	if err != nil{
		glog.V(1).Info(err.Error())
	}
	valEthBTC,err = strconv.ParseFloat(formattedTickers["eth_btc"].AskPrice,64)
	if err != nil{
		glog.V(1).Info(err.Error())
	}

	for index := 0; index < len(symbols); index++ {
		symbol := symbols[index]
		BtcEthBtcArbitrage(formattedTickers,Infos,symbol)
		UsdtToBtcEthToUsdt(formattedTickers,Infos,symbol,"eth")
		UsdtToBtcEthToUsdt(formattedTickers,Infos,symbol,"btc")
		EthBtcToUsdtBtcToEthBtc(formattedTickers,Infos,symbol,"eth","btc")
		EthBtcToUsdtBtcToEthBtc(formattedTickers,Infos,symbol,"btc","usdt")
		EthBtcToUsdtBtcToEthBtc(formattedTickers,Infos,symbol,"eth","usdt")
	}

}