package gocryptobot_tests

import (
	"encoding/json"
	"fmt"
	"github.com/florianpautot/http"
	"github.com/florianpautot/model/requests"
	"github.com/florianpautot/utils"
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
	fmt.Println(res)
}

func TestHTTPGetSecured(t *testing.T) {
	config, err:= utils.LoadConfig("../config.yaml")
	http.Config = config
	now := strconv.FormatInt(time.Now().Unix()*1000,10)

	res,err := http.HTTPGet(config.APIEndpoint+"/api/v1/account","?ts="+now,true)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(res)
}

func TestHTTPPostSecured(t *testing.T) {
	config, err:= utils.LoadConfig("../config.yaml")
	http.Config = config

	order := requests.Order{
		Symbol:    "eth_btc",
		Side:      "sell",
		Type:      "limit",
		Price:     99999999999,
		Quantity:  0.01,
		Timestamp: strconv.FormatInt(time.Now().Unix()*1000,10),
	}

	data,err := json.Marshal(order)
	if err != nil {
		t.Error(err.Error())
	}
	res,err := http.HTTPPost(config.APIEndpoint+"/api/v1/order",data)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(res)
}

func TestHTTPDelete(t *testing.T) {
	config, err:= utils.LoadConfig("../config.yaml")
	http.Config = config
	now := strconv.FormatInt(time.Now().Unix()*1000,10)

	res,err := http.HTTPDelete(config.APIEndpoint+"/api/v1/account","?ts="+now)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(res)
}
