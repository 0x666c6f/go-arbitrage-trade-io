package arbitrage

import (
	"fmt"
	"github.com/florianpautot/go-arbitrage-trade-io/model"
	"github.com/florianpautot/go-arbitrage-trade-io/model/responses"
	"github.com/florianpautot/go-arbitrage-trade-io/tradeio"
	"github.com/florianpautot/go-arbitrage-trade-io/utils"
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
	restartDate := time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),time.Now().Hour(),time.Now().Minute(),Config.StartSecond,0,time.UTC)
	fmt.Println("Starting arbitrage")

	for TotalMinuteWeight < (Config.APIMinuteLimit - 23) && Config.EndSecond > time.Now().Second() {
		launchArbitrages()
		if Config.Timeout != ""{
			duration,err:= time.ParseDuration(Config.Timeout)
			if err != nil {
				fmt.Errorf(err.Error())
			}
			time.Sleep(duration)
		}
	}

	balances, err := tradeio.Account()
	if err != nil {
		fmt.Errorf(err.Error())
	}

	if len(balances.Balances) > 0 {
		formattedBalances := utils.FormatBalance(balances.Balances)
		Config.MaxBTC = formattedBalances["btc"].Available;
		Config.MaxUSDT = formattedBalances["usdt"].Available;
		Config.MaxETH = formattedBalances["eth"].Available;
	}

	TotalMinuteWeight = 0;
	TotalMinuteOrderWeight = 0;
	if time.Now().Second() < restartDate.Second() {
		sleepTime := restartDate.Sub(time.Now())
		fmt.Println("Will sleep", sleepTime.Seconds(), "to reset minute weight");
		time.Sleep(sleepTime)
		fmt.Println("Waking up, sleep is over !");
	}
	Start()
}

func launchArbitrages(){
	tickers,err := tradeio.Tickers()
	if err != nil{
		fmt.Errorf(err.Error())
	}
	TotalMinuteWeight +=20
	if tickers.Code != 0 {
		fmt.Println("Error while retrieving tickers, will sleep until next loop")
		wakeUp := time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),time.Now().Hour(),time.Now().Minute(),Config.StartSecond+1,0,time.UTC)
		fmt.Println("Will sleep",wakeUp)
		time.Sleep(wakeUp.Sub(time.Now()))
		fmt.Println("Waking up, back to work !")
		return
	}

	formattedTickers,symbols := utils.FormatTickers(tickers.Tickers)
	valBTC = formattedTickers["btc_usdt"].AskPrice
	valETH = formattedTickers["eth_usdt"].AskPrice
	valEthBTC = formattedTickers["eth_btc"].AskPrice

	for index := 0; index < len(symbols); index++ {
		symbol := symbols[index]
		BtcEthBtcArbitrage(formattedTickers,Infos,symbol)
		//	await manageArbitrageUSDT_X_Intermediate_USDT(formattedTickers, infos, ticker, "btc");
		//	await manageArbitrageBTCtoXtoETHtoBTC(formattedTickers, infos, ticker);
		//	await manageArbitrageSource_X_Intermediate_Source(formattedTickers, infos, ticker, "eth", "btc");
		//	await manageArbitrageSource_X_Intermediate_Source(formattedTickers, infos, ticker, "btc", "usdt");
		//	await manageArbitrageSource_X_Intermediate_Source(formattedTickers, infos, ticker, "eth", "usdt");
		//	await manageArbitrageUSDT_X_Intermediate_USDT(formattedTickers, infos, ticker, "eth");
	}

}