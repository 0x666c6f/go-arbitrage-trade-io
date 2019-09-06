package arbitrage

import (
	"fmt"
	"github.com/florianpautot/go-arbitrage-trade-io/model/requests"
	"github.com/florianpautot/go-arbitrage-trade-io/model/responses"
	"github.com/florianpautot/go-arbitrage-trade-io/tradeio"
	"github.com/florianpautot/go-arbitrage-trade-io/utils"
	"sort"
	"time"
)


func BtcEthBtcArbitrage(tickers map[string]responses.Ticker, infos map[string]responses.Symbol, symbol string) {
	tickerBTC := tickers[symbol+"_btc"]
	tickerETH := tickers[symbol+"_eth"]
	tickerEthBtc := tickers["eth_btc"]

	if tickerETH != (responses.Ticker{}) &&
		tickerBTC != (responses.Ticker{}) &&
		tickerEthBtc != (responses.Ticker{}) &&
		tickerEthBtc.BidPrice > 0 &&
		tickerBTC.AskPrice > 0 &&
		tickerETH.BidPrice > 0 {
		precBTC := infos[symbol+"_btc"].BaseAssetPrecision
		precETH := infos[symbol+"_eth"].BaseAssetPrecision
		precETHBTC := infos["eth_btc"].BaseAssetPrecision

		bonus := tickerETH.BidPrice * tickerEthBtc.BidPrice / tickerBTC.AskPrice

		if bonus > Config.MinProfit {
			if tickerBTC.AskPrice*tickerBTC.AskQty > Config.MinBTC &&
				tickerETH.BidQty*tickerETH.BidQty > Config.MinETH &&
				tickerETH.BidPrice*tickerETH.BidQty*valEthBTC > Config.MinBTC {

				price := tickerBTC.AskPrice
				mins := []float64{Config.MaxBTC / price, tickerBTC.AskQty, tickerETH.BidQty}
				sort.Float64s(mins)
				qty := utils.RoundUp(utils.RoundDown(mins[0], precBTC), precETH)

				TotalMinuteWeight++
				TotalMinuteOrderWeight++

				orderA := requests.Order{
					Symbol:    symbol + "_btc",
					Side:      "buy",
					Type:      "limit",
					Price:     price,
					Quantity:  qty,
					Timestamp: time.Now().Unix() * 1000,
				}
				orderAResp, err := tradeio.Order(orderA)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				if orderAResp.Code == 0 && orderAResp.Order.Status == "Completed" {
					price = tickerETH.BidPrice
					qty = utils.RoundDown(orderAResp.Order.BaseAmount-orderAResp.Order.Commission, precETH)

					TotalMinuteWeight++
					TotalMinuteOrderWeight++

					orderB := requests.Order{
						Symbol:    symbol + "_eth",
						Side:      "sell",
						Type:      "limit",
						Price:     price,
						Quantity:  qty,
						Timestamp: time.Now().Unix() * 1000,
					}
					orderBResp, err := tradeio.Order(orderB)
					if err != nil {
						fmt.Println(err.Error())
						return
					}

					if orderBResp.Code == 0 && orderBResp.Order.Status == "Completed" {
						price = tickerEthBtc.BidPrice
						qty = utils.RoundUp(orderBResp.Order.Total-orderBResp.Order.Commission, precETHBTC)

						TotalMinuteWeight++
						TotalMinuteOrderWeight++

						orderC := requests.Order{
							Symbol:    "eth_btc",
							Side:      "sell",
							Type:      "limit",
							Price:     price,
							Quantity:  qty,
							Timestamp: time.Now().Unix() * 1000,
						}
						_, err := tradeio.Order(orderC)
						if err != nil {
							fmt.Println(err.Error())
							return
						}
					}
				} else {
					if orderAResp.Code == 0 && orderAResp.Order.Status == "Working" && orderAResp.Order.UnitsFilled <= 0{
						TotalMinuteWeight++
						TotalMinuteOrderWeight++
						_, err := tradeio.CancelOrder(orderAResp.Order.OrderID)
						if err != nil {
							fmt.Println(err.Error())
						}
					}
				}
			}
		}

	}
}
