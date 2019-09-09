package arbitrage

import (
	"github.com/florianpautot/go-arbitrage-trade-io/model/requests"
	"github.com/florianpautot/go-arbitrage-trade-io/model/responses"
	"github.com/florianpautot/go-arbitrage-trade-io/tradeio"
	"github.com/florianpautot/go-arbitrage-trade-io/utils"
	"github.com/golang/glog"
	"sort"
	"strconv"
	"time"
)


func BtcEthBtcArbitrage(tickers map[string]responses.Ticker, infos map[string]responses.Symbol, symbol string) {
	tickerBTC := tickers[symbol+"_btc"]
	tickerETH := tickers[symbol+"_eth"]
	tickerEthBtc := tickers["eth_btc"]

	var orderAResp responses.OrderResponse
	var orderBResp responses.OrderResponse
	var orderCResp responses.OrderResponse

	if tickerETH != (responses.Ticker{}) &&
		tickerBTC != (responses.Ticker{}) &&
		tickerEthBtc != (responses.Ticker{}) {
		precBTC := infos[symbol+"_btc"].BaseAssetPrecision
		precETH := infos[symbol+"_eth"].BaseAssetPrecision
		precETHBTC := infos["eth_btc"].BaseAssetPrecision

		askBtc,err := strconv.ParseFloat(tickerBTC.AskPrice,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		askBtcQty,err := strconv.ParseFloat(tickerBTC.AskQty,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidEth,err := strconv.ParseFloat(tickerETH.BidPrice,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidEthQty,err := strconv.ParseFloat(tickerETH.BidQty,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidEthBtc, err := strconv.ParseFloat(tickerEthBtc.BidPrice,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		if bidEthBtc > 0 &&
			askBtc > 0 &&
			bidEth > 0{
			bonus := bidEth * bidEthBtc / askBtc
			glog.V(3).Info(symbol, " Bonus = ", bonus)

			if bonus > Config.MinProfit {
				if askBtc*askBtcQty > Config.MinBTC &&
					bidEth*bidEthQty > Config.MinETH &&
					bidEth*bidEthQty*valEthBTC > Config.MinBTC {

					price := askBtc
					mins := []float64{Config.MaxBTC / price, askBtcQty, bidEthQty}
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
					orderAResp, err = tradeio.Order(orderA)
					if err != nil {
						glog.V(2).Info(err.Error())
						return
					}
					glog.V(2).Info(symbol, " Order A = ", orderAResp)

					if orderAResp.Code == 0 && orderAResp.Order.Status == "Completed" {
						price = bidEth
						orderAAmount,err := strconv.ParseFloat(orderAResp.Order.BaseAmount,64)
						if err != nil {
							glog.V(2).Info(err.Error())
							return
						}
						orderACommission, err := strconv.ParseFloat(orderAResp.Order.Commission,64)
						if err != nil {
							glog.V(2).Info(err.Error())
							return
						}
						qty = utils.RoundDown(orderAAmount-orderACommission, precETH)

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
						orderBResp, err = tradeio.Order(orderB)
						if err != nil {
							glog.V(2).Info(err.Error())
							return
						}
						glog.V(2).Info(symbol, " Order B = ", orderBResp)

						if orderBResp.Code == 0 && orderBResp.Order.Status == "Completed" {
							orderBAmount,err := strconv.ParseFloat(orderBResp.Order.Total,64)
							if err != nil {
								glog.V(2).Info(err.Error())
								return
							}
							orderBCommission, err := strconv.ParseFloat(orderAResp.Order.Commission,64)
							if err != nil {
								glog.V(2).Info(err.Error())
								return
							}
							price = bidEthBtc
							qty = utils.RoundUp(orderBAmount-orderBCommission, precETHBTC)

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
							orderCResp, err = tradeio.Order(orderC)
							if err != nil {
								glog.V(2).Info(err.Error())
								return
							}
							glog.V(2).Info(symbol, " Order C = ", orderCResp)

							glog.V(2).Info(symbol, "Arbitrage result : <", symbol,">", " bonus = ", bonus )

						}
					} else {
						orderAfilled,err := strconv.ParseFloat(orderAResp.Order.UnitsFilled,64)
						if err != nil {
							glog.V(2).Info(err.Error())
							return
						}
						if orderAResp.Code == 0 && orderAResp.Order.Status == "Working" && orderAfilled <= 0{
							TotalMinuteWeight++
							TotalMinuteOrderWeight++
							_, err := tradeio.CancelOrder(orderAResp.Order.OrderID)
							if err != nil {
								glog.V(2).Infoln(err.Error())
							}
						}
					}
				} else {
					glog.V(3).Info(symbol, " Quantity is not enough")
				}
			}
		}


	}
}
