package gocryptobot_tests

import (
	"github.com/florianpautot/model/requests"
	"github.com/florianpautot/tradeio"
	"github.com/florianpautot/utils"
	"strconv"
	"testing"
	"time"
)

//TestInfo :
func TestInfo(t *testing.T) {
	config, err := utils.LoadConfig("../config.yaml")
	if err != nil {
		t.Error("Error while getting config:",err)
	}
	tradeio.Config = config
	infos, err := tradeio.Info();
	if err != nil {
		t.Error("Error while getting info:",err)
	}

	if infos.Code != 0{
		t.Error("Error while getting info, expecting code = 0 but got ",infos.Code)
	}
}

func TestTickers(t *testing.T) {
	config, err := utils.LoadConfig("../config.yaml")
	if err != nil {
		t.Error("Error while getting config:",err)
	}
	tradeio.Config = config
	infos, err := tradeio.Tickers();
	if err != nil {
		t.Error("Error while getting tickers:",err)
	}

	if infos.Code != 0{
		t.Error("Error while getting tickers, expecting code = 0 but got ",infos.Code)
	}
}

func TestOrder(t *testing.T) {
	config, err := utils.LoadConfig("../config.yaml")
	if err != nil {
		t.Error("Error while getting config:",err)
	}
	tradeio.Config = config

	order := requests.Order{
		Symbol:    "eth_btc",
		Side:      "sell",
		Type:      "limit",
		Price:     99999999999,
		Quantity:  0.01,
		Timestamp: strconv.FormatInt(time.Now().Unix()*1000,10),
	}

	infos, err := tradeio.Order(order);
	if err != nil {
		t.Error("Error while getting tickers:",err)
	}

	if infos.Code != 0{
		t.Error("Error while getting tickers, expecting code = 0 but got ",infos.Code)
	}
}

func TestCancel(t *testing.T) {

	config, err := utils.LoadConfig("../config.yaml")
	if err != nil {
		t.Error("Error while getting config:",err)
	}
	tradeio.Config = config

	order := requests.Order{
		Symbol:    "eth_btc",
		Side:      "sell",
		Type:      "limit",
		Price:     99999999999,
		Quantity:  0.01,
		Timestamp: strconv.FormatInt(time.Now().Unix()*1000,10),
	}

	infos, err := tradeio.Order(order);
	if err != nil {
		t.Error("Error while getting tickers:",err)
	}

	if infos.Code != 0{
		t.Error("Error while getting tickers, expecting code = 0 but got ",infos.Code)
	}


	cancelResp, err := tradeio.Can();
	if err != nil {
		t.Error("Error while getting tickers:",err)
	}

	if infos.Code != 0{
		t.Error("Error while getting tickers, expecting code = 0 but got ",infos.Code)
	}
}