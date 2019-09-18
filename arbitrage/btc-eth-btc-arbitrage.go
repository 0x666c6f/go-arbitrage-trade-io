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


func BtcEthBtcArbitrage(tickers map[string]binance.BookTicker, infos map[string]binance.Symbol, symbol string) {
	tickerBTC := tickers[symbol+"BTC"]
	tickerETH := tickers[symbol+"ETH"]
	tickerEthBtc := tickers["ETHBTC"]

	var orderA *binance.CreateOrderResponse
	var orderB *binance.CreateOrderResponse
	var orderC *binance.CreateOrderResponse

	if tickerETH != (binance.BookTicker{}) &&
		tickerBTC != (binance.BookTicker{}) &&
		tickerEthBtc != (binance.BookTicker{}) {
		precBTC := infos[symbol+"_btc"].BaseAssetPrecision
		precETH := infos[symbol+"_eth"].BaseAssetPrecision
		precETHBTC := infos["eth_btc"].BaseAssetPrecision

		askBtc,err := strconv.ParseFloat(tickerBTC.AskPrice,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		askBtcQty,err := strconv.ParseFloat(tickerBTC.AskQuantity,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidEth,err := strconv.ParseFloat(tickerETH.BidPrice,64)
		if err != nil {
			glog.V(2).Info(err.Error())
			return
		}
		bidEthQty,err := strconv.ParseFloat(tickerETH.BidQuantity,64)
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

			price := askBtc

			if bonus > global.GlobalConfig.MinProfit {
				if askBtc*askBtcQty > global.GlobalConfig.MinBTC &&
					bidEth*bidEthQty > global.GlobalConfig.MinETH &&
					bidEth*bidEthQty*valEthBTC > global.GlobalConfig.MinBTC &&
					global.GlobalConfig.MaxBTC / price > global.GlobalConfig.MinBTC {

					mins := []float64{utils.RoundDown(global.GlobalConfig.MaxBTC / price, precBTC), askBtcQty, bidEthQty}
					sort.Float64s(mins)
					qty := utils.RoundUp(utils.RoundDown(mins[0], precBTC), precETH)

					if(qty == 0){
						return
					}

					TotalMinuteWeight++
					TotalMinuteOrderWeight++

					orderA, err = global.Binance.NewCreateOrderService().Symbol(symbol+"BTC").
						Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
						TimeInForce(binance.TimeInForceTypeGTC).Quantity(fmt.Sprintf("%f", qty)).
						Price(fmt.Sprintf("%f", price)).Do(context.Background())
					if err != nil {
						glog.V(2).Info(err.Error())
						return
					}

					glog.V(2).Info(symbol, " Order A = ", orderA)

					if orderA.OrderID != 0 && orderA.Status == binance.OrderStatusTypeFilled {
						price = bidEth
						orderAAmount,err := strconv.ParseFloat(orderA.ExecutedQuantity,64)
						if err != nil {
							glog.V(2).Info(err.Error())
							return
						}

						qty = utils.RoundDown(orderAAmount*0.999, precETH)

						TotalMinuteWeight++
						TotalMinuteOrderWeight++

						orderB, err = global.Binance.NewCreateOrderService().Symbol(symbol+"ETH").
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

							price = bidEthBtc
							qty = utils.RoundUp((orderBAmount*0.999)*price, precETHBTC)

							TotalMinuteWeight++
							TotalMinuteOrderWeight++

							orderC, err = global.Binance.NewCreateOrderService().Symbol("ETHBTC").
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
