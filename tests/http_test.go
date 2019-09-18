package gocryptobot_tests

import (
	"encoding/json"
	"github.com/florianpautot/go-arbitrage/http"
	"github.com/florianpautot/go-arbitrage/global"
	"github.com/florianpautot/go-arbitrage/global/tradeio/requests"
	"github.com/florianpautot/go-arbitrage/utils"
	"github.com/golang/glog"
	"strconv"
	"testing"
	"time"
)
func TestHttp(t *testing.T){
	config, err:= utils.LoadConfig("../config.yaml")
	if err != nil {
		return
	}
	global.GlobalConfig = config
	testHTTPGetUnsecured(t)
	testHTTPGetSecured(t)
	testHTTPPostSecured(t)
	testHTTPDelete(t)
}

//TestGenerateSingature :
func testHTTPGetUnsecured(t *testing.T) {
	res,err := http.HTTPGet(global.GlobalConfig.APIEndpoint+"/api/v1/info","",false)
	if err != nil {
		t.Error(err.Error())
	}
	glog.V(2).Info(res)
}

func testHTTPGetSecured(t *testing.T) {
	now := strconv.FormatInt(time.Now().Unix()*1000,10)

	res,err := http.HTTPGet(global.GlobalConfig.APIEndpoint+"/api/v1/account","?ts="+now,true)
	if err != nil {
		t.Error(err.Error())
	}
	glog.V(2).Info(res)
}

func testHTTPPostSecured(t *testing.T) {
	order := requests.Order{
		Symbol:    "eth_btc",
		Side:      "sell",
		Type:      "",
		Price:     99999999999,
		Quantity:  0.01,
		Timestamp: time.Now().Unix()*1000,
	}

	data,err := json.Marshal(&order)
	if err != nil {
		t.Error(err.Error())
	}
	res,err := http.HTTPPost(global.GlobalConfig.APIEndpoint+"/api/v1/order",data)
	if err != nil {
		t.Error(err.Error())
	}
	glog.V(2).Info(res)
}

func testHTTPDelete(t *testing.T) {
	now := strconv.FormatInt(time.Now().Unix()*1000,10)

	res,err := http.HTTPDelete(global.GlobalConfig.APIEndpoint+"/api/v1/account","?ts="+now)
	if err != nil {
		t.Error(err.Error())
	}
	glog.V(2).Info(res)
}
