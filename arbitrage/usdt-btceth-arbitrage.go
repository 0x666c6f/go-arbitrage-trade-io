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

func UsdtToBtcEthToUsdt(tickers map[string]binance.BookTicker, infos map[string]binance.Symbol, symbol string, intermediate string) {
	tickerUSDT := tickers[symbol+"USDT"]
	tickerIntermediate := tickers[symbol+intermediate]
	tickerIntermediateUSDT := tickers[intermediate+"USDT"]

	var orderA *binance.CreateOrderResponse
	var orderB *binance.CreateOrderResponse
	var orderC *binance.CreateOrderResponse

	if tickerUSDT != (binance.BookTicker{}) &&
		tickerIntermediate != (binance.BookTicker{}) &&
		tickerIntermediateUSDT != (binance.BookTicker{}) {

		precUSDT := infos[symbol+"USDT"].BaseAssetPrecision
		precIntermediate := infos[symbol+intermediate].BaseAssetPrecision
		precIntermediateUSDT := infos[intermediate+"USDT"].BaseAssetPrecision

		askUSDT,err := strconv.ParseFloat(tickerUSDT.AskPrice,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		askUSDTQty,err := strconv.ParseFloat(tickerUSDT.AskQuantity,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidIntermediate,err := strconv.ParseFloat(tickerIntermediate.BidPrice,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidIntermediateQty,err := strconv.ParseFloat(tickerIntermediate.BidQuantity,64)
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


			if bonus > global.GlobalConfig.MinProfit {

				var minIntermediate float64
				var valIntermediate float64

				if intermediate == "eth" {
					minIntermediate = global.GlobalConfig.MinETH
					valIntermediate = valETH
				} else {
					minIntermediate = global.GlobalConfig.MinBTC
					valIntermediate = valBTC
				}

				price := askUSDT

				if askUSDT*askUSDTQty > global.GlobalConfig.MinUSDT &&
					bidIntermediate*bidIntermediateQty > minIntermediate &&
					bidIntermediate*bidIntermediateQty*valIntermediate > global.GlobalConfig.MinUSDT &&
					global.GlobalConfig.MaxUSDT/ price >  global.GlobalConfig.MinUSDT {

					mins := []float64{utils.RoundDown(global.GlobalConfig.MaxUSDT/ price, precUSDT), askUSDTQty, bidIntermediateQty}
					sort.Float64s(mins)
					qty := utils.RoundUp(utils.RoundDown(mins[0], precUSDT), precIntermediate)

					if(qty == 0){
						return
					}

					TotalMinuteWeight++
					TotalMinuteOrderWeight++

					orderA, err = global.Binance.NewCreateOrderService().Symbol(symbol+"USDT").
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
						orderAAmount,err := strconv.ParseFloat(orderA.ExecutedQuantity,64)
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

							orderBAmount,err := strconv.ParseFloat(orderB.ExecutedQuantity,64)
							if err != nil {
								glog.V(2).Info(err.Error())
								return
							}

							price = bidIntermediateUSDT
							qty = utils.RoundUp((orderBAmount*0.999)*price, precIntermediateUSDT)

							TotalMinuteWeight++
							TotalMinuteOrderWeight++


							orderC, err = global.Binance.NewCreateOrderService().Symbol(intermediate+"USDT").
								Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).
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
						if orderC.ExecutedQuantity != ""{
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