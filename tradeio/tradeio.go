package tradeio

import (
	"encoding/json"
	"github.com/florianpautot/http"
	"github.com/florianpautot/model"
	"github.com/florianpautot/model/requests"
	"github.com/florianpautot/model/responses"
	"strconv"
	"time"
)

var Config model.Config

//Info :
func  Info() (responses.Infos, error){
	var infos responses.Infos
	http.Config = Config
	res, err := http.HTTPGet(Config.APIEndpoint+"/api/v1/info","",false)
	if err != nil {
		return responses.Infos{}, err
	}

	err = json.Unmarshal(res, &infos)
	if err != nil {
		return responses.Infos{}, err
	}

	return infos, nil
}


//Tickers :
func Tickers() (responses.Tickers, error){
	var tickers responses.Tickers
	http.Config = Config
	res, err := http.HTTPGet(Config.APIEndpoint+"/api/v1/tickers","",false)
	if err != nil {
		return responses.Tickers{}, err
	}

	err = json.Unmarshal(res, &tickers)
	if err != nil {
		return responses.Tickers{}, err
	}

	return tickers, nil
}


//Order :
func Order(order requests.Order) (responses.OrderResponse, error){
	var orderResponse responses.OrderResponse
	http.Config = Config
	marshOrder, err := json.Marshal(order)
	res, err := http.HTTPPost(Config.APIEndpoint+"/api/v1/order",marshOrder)
	if err != nil {
		return responses.OrderResponse{}, err
	}

	err = json.Unmarshal(res, &orderResponse)
	if err != nil {
		return responses.OrderResponse{}, err
	}

	return orderResponse, nil
}

//DeleteOrder :
func CancelOrder(orderID string) (responses.CancelResponse,error){
	var cancelResp responses.CancelResponse
	http.Config = Config

	res, err := http.HTTPDelete(Config.APIEndpoint+"/api/v1/order/"+orderID,"?ts="+strconv.FormatInt(time.Now().Unix()*1000,10))
	if err != nil {
		return responses.CancelResponse{}, err
	}

	err = json.Unmarshal(res, &cancelResp)
	if err != nil {
		return responses.CancelResponse{}, err
	}

	return cancelResp, nil
}


//ClosedOrders :
func ClosedOrders(symbol string, start int, end int, page int, perPage int) {
}

//Account :
func Account() {

}

