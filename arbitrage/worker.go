package arbitrage

import (
	"github.com/adshao/go-binance"
	"github.com/golang/glog"
)

func Workers(symbol string, formattedTickers map[string]binance.BookTicker, infos map[string]binance.Symbol,  finish func()){
	glog.V(3).Info("Launching async worker")
	ack := make(chan bool)
	phase1Chan := make(chan bool)
	phase2Chan := make(chan bool)
	phase3Chan := make(chan bool)

	//Wait for all routine to finish
	go func() {
		glog.V(3).Info("Waiting for workers to end ...")
		<-ack
		glog.V(3).Info("Work finished !")

		close(phase1Chan)
		close(phase2Chan)
		close(phase3Chan)
		close(ack)
		finish()
	}()

	glog.V(3).Info("Arbitrages wave 1 started")
	//Arbitrages wave 1
	// usdt->XXX->btc->usdt
	// btc->XXX->eth->btc
	go func() {
		glog.V(3).Info("async usdt->XXX->btc->usdt")
		UsdtToBtcEthToUsdt(formattedTickers,infos,symbol,"btc")

		phase1Chan <- true
	}()
	go func() {
		glog.V(3).Info("async btc->XXX->eth->btc")
		BtcEthBtcArbitrage(formattedTickers,infos,symbol)
		phase1Chan <- true
	}()
	for i := 0; i < 2; i++{
		<-phase1Chan
		glog.V(3).Info(i, " Wave 1 end")
	}

	//Arbitrages wave 2
	// eth->XXX->btc->eth
	// btc->XXX->usdt->btc
	glog.V(3).Info("Arbitrages wave 2 started")

	go func() {
		glog.V(3).Info("async eth->XXX->btc->eth")
		EthBtcToUsdtBtcToEthBtc(formattedTickers,infos,symbol,"eth","btc")
		phase2Chan <- true
	}()
	go func() {
		glog.V(3).Info("async btc->XXX->usdt->btc")
		EthBtcToUsdtBtcToEthBtc(formattedTickers,infos,symbol,"btc","usdt")
		phase2Chan <- true
	}()
	for i := 0; i < 2; i++{
		<-phase2Chan
		glog.V(3).Info(i, " Wave 2 end")
	}

	//Arbitrages wave 3
	// eth->XXX->usdt->eth
	// usdt->xxx->eth->usdt
	glog.V(3).Info("Arbitrages wave 3 started")

	go func() {
		glog.V(3).Info("async eth->XXX->usdt->eth")
		EthBtcToUsdtBtcToEthBtc(formattedTickers,infos,symbol,"eth","usdt")

		phase3Chan <- true
	}()
	go func() {
		glog.V(3).Info("async usdt->xxx->eth->usdt")
		UsdtToBtcEthToUsdt(formattedTickers,infos,symbol,"eth")
		phase3Chan <- true
	}()
	for i := 0; i < 2; i++{
		<-phase3Chan
		glog.V(3).Info(i, " Wave 3 end")
	}
	glog.V(3).Info("All waves are finished")
	ack <- true



	return
}