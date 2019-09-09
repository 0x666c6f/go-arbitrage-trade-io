package gocryptobot_tests

import (
	"encoding/json"
	"github.com/florianpautot/go-arbitrage-trade-io/http"
	"github.com/florianpautot/go-arbitrage-trade-io/model/requests"
	"github.com/florianpautot/go-arbitrage-trade-io/utils"
	"github.com/golang/glog"
	"strconv"
	"testing"
	"time"
)

//TestGenerateSingature :
func TestHTTPGetUnsecured(t *testing.T) {
	config, err:= utils.LoadConfig("../config.yaml")
	http.Config = config
	res,err := http.HTTPGet(config.APIEndpoint+"/api/v1/info","",false)
	if err != nil {
		t.Error(err.Error())
	}
	glog.V(2).Info(res)
}

func TestHTTPGetSecured(t *testing.T) {
	config, err:= utils.LoadConfig("../config.yaml")
	http.Config = config
	now := strconv.FormatInt(time.Now().Unix()*1000,10)

	res,err := http.HTTPGet(config.APIEndpoint+"/api/v1/account","?ts="+now,true)
	if err != nil {
		t.Error(err.Error())
	}
	glog.V(2).Info(res)
}

func TestHTTPPostSecured(t *testing.T) {
	config, err:= utils.LoadConfig("../config.yaml")
	http.Config = config

	order := requests.Order{
		Symbol:    "eth_btc",
		Side:      "sell",
		Type:      "",
		Price:     99999999999,
		Quantity:  0.01,
		Timestamp: time.Now().Unix()*1000,
	}

	data,err := json.Marshal(order)
	if err != nil {
		t.Error(err.Error())
	}
	res,err := http.HTTPPost(config.APIEndpoint+"/api/v1/order",data)
	if err != nil {
		t.Error(err.Error())
	}
	glog.V(2).Info(res)
}

func TestHTTPDelete(t *testing.T) {
	config, err:= utils.LoadConfig("../config.yaml")
	http.Config = config
	now := strconv.FormatInt(time.Now().Unix()*1000,10)

	res,err := http.HTTPDelete(config.APIEndpoint+"/api/v1/account","?ts="+now)
	if err != nil {
		t.Error(err.Error())
	}
	glog.V(2).Info(res)
}
