package gocryptobot_tests

import (
	"github.com/florianpautot/go-arbitrage/global"
	"github.com/florianpautot/go-arbitrage/utils"
	"testing"
	"time"
)

func TestTradeIO(t *testing.T){
	config, err:= utils.LoadConfig("../config.yaml")
	if err != nil {
		return
	}
	global.GlobalConfig = config
	TestInfo(t)
	TestTickers(t)
	TestOrder(t)
	TestFailedOrder(t)
	TestCancel(t)
	TestBalances(t)
}

//TestInfo :
func TestInfo(t *testing.T) {
	infos, err := tradeio.Info();
	if err != nil {
		t.Error("Error while getting info:",err)
	}

	if infos.Code != 0{
		t.Error("Error while getting info, expecting code = 0 but got ",infos.Code)
	}
}

func TestTickers(t *testing.T) {
	infos, err := tradeio.Tickers();
	if err != nil {
		t.Error("Error while getting tickers:",err)
	}

	if infos.Code != 0{
		t.Error("Error while getting tickers, expecting code = 0 but got ",infos.Code)
	}
}

func TestOrder(t *testing.T) {
	order := requests.Order{
		Symbol:    "eth_btc",
		Side:      "sell",
		Type:      "limit",
		Price:     99999999999,
		Quantity:  0.01,
		Timestamp:time.Now().Unix()*1000,
	}

	orderResp, err := tradeio.Order(order);
	if err != nil {
		t.Error("Error while creating order:",err)
	}

	if orderResp.Code != 0{
		t.Error("Error while creating order, expecting code = 0 but got ",orderResp.Code)
	}
}

func TestFailedOrder(t *testing.T) {
	order := requests.Order{
		Symbol:    "eth_btc",
		Side:      "sell",
		Type:      "limit",
		Price:     99999999999,
		Quantity:  99999999999,
		Timestamp:time.Now().Unix()*1000,
	}

	orderResp, err := tradeio.Order(order);
	if err != nil {
		t.Error("Error while creating order:",err)
	}  else if orderResp != (responses.OrderResponse{}){
		t.Error("Exptected an error but got empty error")
		t.Fail()
	}
}

func TestCancel(t *testing.T) {
	order := requests.Order{
		Symbol:    "eth_btc",
		Side:      "sell",
		Type:      "limit",
		Price:     99999999999,
		Quantity:  0.01,
		Timestamp: time.Now().Unix()*1000,
	}

	orderResp, err := tradeio.Order(order);
	if err != nil {
		t.Error("Error while creating order to cancel:",err)
	}
	if orderResp.Code != 0{
		t.Error("Error while getting creating order to cancel, expecting code = 0 but got ",orderResp.Code)
	}

	cancelResp, err := tradeio.CancelOrder(orderResp.Order.OrderID);
	if err != nil {
		t.Error("Error while cancelling order:",err)
	}
	if cancelResp.Code != 0{
		t.Error("Error while cancelling order, expecting code = 0 but got ",cancelResp.Code)
	}
}

func TestBalances(t *testing.T) {
	balances, err := tradeio.Account();
	if err != nil {
		t.Error("Error while getting balances:",err)
	}
	formattedBalances := utils.FormatBalance(balances.Balances)

	if len(formattedBalances) == 0{
		t.Error("Error while getting balances, expecting len > 0 but got ",len(formattedBalances))
	}
}