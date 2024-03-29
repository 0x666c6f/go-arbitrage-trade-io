package arbitrage

import (
	"github.com/florianpautot/go-arbitrage-trade-io/model"
	"github.com/florianpautot/go-arbitrage-trade-io/model/requests"
	"github.com/florianpautot/go-arbitrage-trade-io/model/responses"
	"github.com/florianpautot/go-arbitrage-trade-io/tradeio"
	"github.com/florianpautot/go-arbitrage-trade-io/utils"
	"github.com/golang/glog"
	"sort"
	"strconv"
	"time"
)

func EthBtcToUsdtBtcToEthBtc(tickers map[string]responses.Ticker, infos map[string]responses.Symbol, symbol string, source string, intermediate string) {
	tickerSource := tickers[symbol+"_"+source]
	tickerIntermediate := tickers[symbol+"_"+intermediate]
	tickerSourceIntermediate := tickers[source+"_"+intermediate]

	var orderAResp responses.OrderResponse
	var orderBResp responses.OrderResponse
	var orderCResp responses.OrderResponse

	if tickerSource != (responses.Ticker{}) &&
		tickerIntermediate != (responses.Ticker{}) &&
		tickerSourceIntermediate != (responses.Ticker{}) {

		precSource := infos[symbol+"_"+source].BaseAssetPrecision
		precIntermediate := infos[symbol+"_"+intermediate].BaseAssetPrecision
		precSourceIntermediate := infos[source+"_"+intermediate].BaseAssetPrecision

		askSource, err := strconv.ParseFloat(tickerSource.AskPrice, 64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		askSourceQty, err := strconv.ParseFloat(tickerSource.AskQty, 64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidIntermediate, err := strconv.ParseFloat(tickerIntermediate.BidPrice, 64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidIntermediateQty, err := strconv.ParseFloat(tickerIntermediate.BidQty, 64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		askSourceIntermediate, err := strconv.ParseFloat(tickerSourceIntermediate.AskPrice, 64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		if bidIntermediate > 0 &&
			askSourceIntermediate > 0 &&
			askSource > 0 {

			bonus := bidIntermediate / askSource / askSourceIntermediate
			glog.V(3).Info(symbol, " Bonus = ", bonus)

			if bonus > model.GlobalConfig.MinProfit {
				var minIntermediate float64

				if intermediate == "eth" {
					minIntermediate = model.GlobalConfig.MinETH
				} else if intermediate == "btc" {
					minIntermediate = model.GlobalConfig.MinBTC
				} else {
					minIntermediate = model.GlobalConfig.MinUSDT
				}

				var minSource float64
				var maxSource float64

				if source == "eth" {
					minSource = model.GlobalConfig.MinETH
					maxSource = model.GlobalConfig.MaxETH
				} else if source == "btc" {
					minSource = model.GlobalConfig.MinBTC
					maxSource = model.GlobalConfig.MaxBTC
				} else {
					minSource = model.GlobalConfig.MinUSDT
					maxSource = model.GlobalConfig.MaxUSDT
				}

				var valSourceIntermediate float64
				if source == "eth" && intermediate == "usdt" {
					valSourceIntermediate = valETH;
				} else if source == "btc" && intermediate == "usdt" {
					valSourceIntermediate = valBTC;
				} else{
					valSourceIntermediate = valEthBTC;
				}

				price := askSource

				if bidIntermediate*bidIntermediateQty > minIntermediate &&
					askSource * askSourceQty > minSource &&
					askSource *askSourceQty * valSourceIntermediate > minIntermediate &&
					maxSource/price > minSource {

					mins := []float64{utils.RoundDown(maxSource / price, precSource), askSourceQty, bidIntermediateQty}
					sort.Float64s(mins)
					qty := utils.RoundUp(utils.RoundDown(mins[0], precSource), precIntermediate)

					if(qty == 0){
						return
					}

					TotalMinuteWeight++
					TotalMinuteOrderWeight++

					orderA := requests.Order{
						Symbol:    symbol + "_"+source,
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
					glog.V(3).Info(symbol, " Order A Resp = ", orderAResp)

					if orderAResp.Code == 0 && orderAResp.Order.Status == "Completed" {

						price = bidIntermediate
						orderAAmount, err := strconv.ParseFloat(orderAResp.Order.BaseAmount, 64)
						if err != nil {
							glog.V(2).Info(err.Error())
							return
						}
						orderACommission, err := strconv.ParseFloat(orderAResp.Order.Commission, 64)
						if err != nil {
							glog.V(2).Info(err.Error())
							return
						}
						qty = utils.RoundDown(orderAAmount-orderACommission, precIntermediate)

						TotalMinuteWeight++
						TotalMinuteOrderWeight++

						orderB := requests.Order{
							Symbol:    symbol + "_" + intermediate,
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
						glog.V(3).Info(symbol, " Order B = ", orderBResp)

						if orderBResp.Code == 0 && orderBResp.Order.Status == "Completed" {

							orderBAmount, err := strconv.ParseFloat(orderBResp.Order.Total, 64)
							if err != nil {
								glog.V(2).Info(err.Error())
								return
							}
							orderBCommission, err := strconv.ParseFloat(orderBResp.Order.Commission, 64)
							if err != nil {
								glog.V(2).Info(err.Error())
								return
							}
							price = askSourceIntermediate
							qty = utils.RoundUp((orderBAmount-orderBCommission)/price, precSourceIntermediate)

							TotalMinuteWeight++
							TotalMinuteOrderWeight++

							orderC := requests.Order{
								Symbol:    source + "_"+intermediate,
								Side:      "buy",
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
							glog.V(3).Info(symbol, " Order C = ", orderCResp)

							glog.V(1).Info("Arbitrage result : <", symbol,">", " bonus = ", bonus )

						}
					} else {
						if orderAResp.Order.UnitsFilled != ""{
							orderAfilled, err := strconv.ParseFloat(orderAResp.Order.UnitsFilled, 64)
							if err != nil {
								glog.V(2).Info(err.Error())
								return
							}
							if orderAResp.Code == 0 && orderAResp.Order.Status == "Working" && orderAfilled <= 0 {
								TotalMinuteWeight++
								TotalMinuteOrderWeight++
								_, err := tradeio.CancelOrder(orderAResp.Order.OrderID)
								if err != nil {
									glog.V(2).Infoln(err.Error())
								}
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
