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

func UsdtToBtcEthToUsdt(tickers map[string]responses.Ticker, infos map[string]responses.Symbol, symbol string, intermediate string) {
	tickerUSDT := tickers[symbol+"_usdt"]
	tickerIntermediate := tickers[symbol+"_"+intermediate]
	tickerIntermediateUSDT := tickers[intermediate+"_usdt"]

	if tickerUSDT != (responses.Ticker{}) &&
		tickerIntermediate != (responses.Ticker{}) &&
		tickerIntermediateUSDT != (responses.Ticker{}) {

		precUSDT := infos[symbol+"_usdt"].BaseAssetPrecision
		precIntermediate := infos[symbol+"_"+intermediate].BaseAssetPrecision
		precIntermediateUSDT := infos[intermediate+"_usdt"].BaseAssetPrecision

		askUSDT,err := strconv.ParseFloat(tickerUSDT.AskPrice,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		askUSDTQty,err := strconv.ParseFloat(tickerUSDT.AskQty,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidIntermediate,err := strconv.ParseFloat(tickerIntermediate.BidPrice,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidIntermediateQty,err := strconv.ParseFloat(tickerIntermediate.BidQty,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidIntermediateUSDT, err := strconv.ParseFloat(tickerIntermediateUSDT.BidPrice,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}

		if bidIntermediate > 0 &&
			bidIntermediateUSDT > 0 &&
			askUSDT > 0{

			bonus := bidIntermediateUSDT * bidIntermediate / askUSDT
			glog.V(3).Info(symbol, " Bonus = ", bonus)


			if bonus > model.GlobalConfig.MinProfit {

				var minIntermediate float64
				var valIntermediate float64

				if intermediate == "eth" {
					minIntermediate = model.GlobalConfig.MinETH
					valIntermediate = valETH
				} else {
					minIntermediate = model.GlobalConfig.MinBTC
					valIntermediate = valBTC
				}

				price := askUSDT

				if askUSDT*askUSDTQty > model.GlobalConfig.MinUSDT &&
					bidIntermediate*bidIntermediateQty > minIntermediate &&
					bidIntermediate*bidIntermediateQty*valIntermediate > model.GlobalConfig.MinUSDT &&
					model.GlobalConfig.MaxUSDT/ price >  model.GlobalConfig.MinUSDT {

					mins := []float64{utils.RoundDown(model.GlobalConfig.MaxUSDT/ price, precUSDT), askUSDTQty, bidIntermediateQty}
					sort.Float64s(mins)
					qty := utils.RoundUp(utils.RoundDown(mins[0], precUSDT), precIntermediate)

					if(qty == 0){
						return
					}

					TotalMinuteWeight++
					TotalMinuteOrderWeight++

					orderA := requests.Order{
						Symbol:    symbol + "_usdt",
						Side:      "buy",
						Type:      "limit",
						Price:     price,
						Quantity:  qty,
						Timestamp: time.Now().Unix() * 1000,
					}

					orderAResp, err := tradeio.Order(orderA)
					if err != nil {
						glog.V(2).Info(err.Error())
						return
					}
					glog.V(3).Info(symbol, " Order A = ", orderAResp)

					if orderAResp.Code == 0 && orderAResp.Order.Status == "Completed" {
						price = bidIntermediate
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
						qty = utils.RoundDown(orderAAmount-orderACommission, precIntermediate)

						TotalMinuteWeight++
						TotalMinuteOrderWeight++

						orderB := requests.Order{
							Symbol:    symbol + "_"+intermediate,
							Side:      "sell",
							Type:      "limit",
							Price:     price,
							Quantity:  qty,
							Timestamp: time.Now().Unix() * 1000,
						}

						orderBResp, err := tradeio.Order(orderB)
						if err != nil {
							glog.V(2).Info(err.Error())
							return
						}
						glog.V(3).Info(symbol, " Order B = ", orderBResp)

						if orderBResp.Code == 0 && orderBResp.Order.Status == "Completed" {

							orderBAmount,err := strconv.ParseFloat(orderBResp.Order.Total,64)
							if err != nil {
								glog.V(2).Info(err.Error())
								return
							}
							orderBCommission, err := strconv.ParseFloat(orderBResp.Order.Commission,64)
							if err != nil {
								glog.V(2).Info(err.Error())
								return
							}
							price = bidIntermediateUSDT
							qty = utils.RoundUp(orderBAmount-orderBCommission, precIntermediateUSDT)

							TotalMinuteWeight++
							TotalMinuteOrderWeight++

							orderC := requests.Order{
								Symbol:    intermediate+"_usdt",
								Side:      "sell",
								Type:      "limit",
								Price:     price,
								Quantity:  qty,
								Timestamp: time.Now().Unix() * 1000,
							}

							orderCResp, err := tradeio.Order(orderC)
							if err != nil {
								glog.V(2).Info(err.Error())
								return
							}
							glog.V(3).Info(symbol, " Order C = ", orderCResp)

							glog.V(1).Info("Successful Arbitrage result : <", symbol,">", " bonus = ", bonus )
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