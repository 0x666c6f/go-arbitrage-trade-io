package arbitrage

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance"
	"github.com/florianpautot/go-arbitrage/global"
	"github.com/florianpautot/go-arbitrage/utils"
	"github.com/golang/glog"
	"sort"
	"strconv"
)

func EthBtcToUsdtBtcToEthBtc(tickers map[string]binance.BookTicker, infos map[string]binance.Symbol, symbol string, source string, intermediate string) {
	tickerSource := tickers[symbol+source]
	tickerIntermediate := tickers[symbol+intermediate]
	tickerSourceIntermediate := tickers[source+intermediate]

	var orderA *binance.CreateOrderResponse
	var orderB *binance.CreateOrderResponse
	var orderC *binance.CreateOrderResponse

	if tickerSource != (binance.BookTicker{}) &&
		tickerIntermediate != (binance.BookTicker{}) &&
		tickerSourceIntermediate != (binance.BookTicker{}) {

		precSource := infos[symbol+source].BaseAssetPrecision
		precIntermediate := infos[symbol+intermediate].BaseAssetPrecision
		precSourceIntermediate := infos[source+"_"+intermediate].BaseAssetPrecision

		askSource, err := strconv.ParseFloat(tickerSource.AskPrice, 64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		askSourceQty, err := strconv.ParseFloat(tickerSource.AskQuantity, 64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidIntermediate, err := strconv.ParseFloat(tickerIntermediate.BidPrice, 64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidIntermediateQty, err := strconv.ParseFloat(tickerIntermediate.BidQuantity, 64)
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

			if bonus > global.GlobalConfig.MinProfit {
				var minIntermediate float64

				if intermediate == "eth" {
					minIntermediate = global.GlobalConfig.MinETH
				} else if intermediate == "btc" {
					minIntermediate = global.GlobalConfig.MinBTC
				} else {
					minIntermediate = global.GlobalConfig.MinUSDT
				}

				var minSource float64
				var maxSource float64

				if source == "eth" {
					minSource = global.GlobalConfig.MinETH
					maxSource = global.GlobalConfig.MaxETH
				} else if source == "btc" {
					minSource = global.GlobalConfig.MinBTC
					maxSource = global.GlobalConfig.MaxBTC
				} else {
					minSource = global.GlobalConfig.MinUSDT
					maxSource = global.GlobalConfig.MaxUSDT
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

					orderA, err = global.Binance.NewCreateOrderService().Symbol(symbol+source).
						Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
						TimeInForce(binance.TimeInForceTypeGTC).Quantity(fmt.Sprintf("%f", qty)).
						Price(fmt.Sprintf("%f", price)).Do(context.Background())
					if err != nil {
						glog.V(2).Info(err.Error())
						return
					}
					glog.V(2).Info(symbol, " Order A = ", orderA)

					if orderA.OrderID != 0 && orderA.Status == binance.OrderStatusTypeFilled {

						price = bidIntermediate
						orderAAmount, err := strconv.ParseFloat(orderA.ExecutedQuantity, 64)
						if err != nil {
							glog.V(2).Info(err.Error())
							return
						}
						qty = utils.RoundDown(orderAAmount*0.999, precIntermediate)

						TotalMinuteWeight++
						TotalMinuteOrderWeight++

						orderB, err = global.Binance.NewCreateOrderService().Symbol(symbol+intermediate).
							Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).
							TimeInForce(binance.TimeInForceTypeGTC).Quantity(fmt.Sprintf("%f", qty)).
							Price(fmt.Sprintf("%f", price)).Do(context.Background())
						if err != nil {
							glog.V(2).Info(err.Error())
							return
						}
						glog.V(2).Info(symbol, " Order B = ", orderB)

						if orderB.OrderID != 0 && orderB.Status == binance.OrderStatusTypeFilled {

							orderBAmount, err := strconv.ParseFloat(orderB.ExecutedQuantity, 64)
							if err != nil {
								glog.V(2).Info(err.Error())
								return
							}

							price = askSourceIntermediate
							qty = utils.RoundUp((orderBAmount*0.999)/price, precSourceIntermediate)

							TotalMinuteWeight++
							TotalMinuteOrderWeight++

							orderC, err = global.Binance.NewCreateOrderService().Symbol(source+intermediate).
								Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
								TimeInForce(binance.TimeInForceTypeGTC).Quantity(fmt.Sprintf("%f", qty)).
								Price(fmt.Sprintf("%f", price)).Do(context.Background())
							if err != nil {
								glog.V(2).Info(err.Error())
								return
							}
							glog.V(2).Info(symbol, " Order C = ", orderC)

							glog.V(1).Info("Successful Arbitrage result : <", symbol,">", " bonus = ", bonus )

						}
					} else {
						if orderA.ExecutedQuantity != ""{
							orderAfilled, err := strconv.ParseFloat(orderA.ExecutedQuantity, 64)
							if err != nil {
								glog.V(2).Info(err.Error())
								return
							}
							if orderA.OrderID == 0 && orderA.Status == binance.OrderStatusTypePartiallyFilled && orderAfilled <= 0 {
								TotalMinuteWeight++
								TotalMinuteOrderWeight++
								_,err := global.Binance.NewCancelOrderService().OrderID(orderA.OrderID).Do(context.Background())
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
