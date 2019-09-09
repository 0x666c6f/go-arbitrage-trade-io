package tradeio

import (
	"encoding/json"
	"errors"
	"github.com/florianpautot/go-arbitrage-trade-io/http"
	"github.com/florianpautot/go-arbitrage-trade-io/model"
	"github.com/florianpautot/go-arbitrage-trade-io/model/requests"
	"github.com/florianpautot/go-arbitrage-trade-io/model/responses"
	"strconv"
	"time"
)

var Config model.Config

//Info :
func  Info() (responses.Infos, error){
	var infos responses.Infos
	var errorResponse responses.ErrorResponse

	http.Config = Config
	res, err := http.HTTPGet(Config.APIEndpoint+"/api/v1/info","",false)
	if err != nil {
		return responses.Infos{}, err
	}

	err = json.Unmarshal(res, &errorResponse)
	if err != nil {
		return responses.Infos{}, err
	}

	if len(errorResponse.Error) > 0 {
		return responses.Infos{}, errors.New(errorResponse.Error)

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
	var errorResponse responses.ErrorResponse

	http.Config = Config
	res, err := http.HTTPGet(Config.APIEndpoint+"/api/v1/tickers","",false)
	if err != nil {
		return responses.Tickers{}, err
	}

	err = json.Unmarshal(res, &errorResponse)
	if err != nil {
		return responses.Tickers{}, err
	}

	if len(errorResponse.Error) > 0 {
		return responses.Tickers{}, errors.New(errorResponse.Error)

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
	var errorResponse responses.ErrorResponse
	http.Config = Config
	marshOrder, err := json.Marshal(order)
	res, err := http.HTTPPost(Config.APIEndpoint+"/api/v1/order",marshOrder)
	if err != nil {
		return responses.OrderResponse{}, err
	}

	err = json.Unmarshal(res, &errorResponse)
	if err != nil {
		return responses.OrderResponse{}, err
	}

	if len(errorResponse.Error) > 0 {
		return responses.OrderResponse{}, errors.New(errorResponse.Error)
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
	var errorResponse responses.ErrorResponse

	http.Config = Config

	res, err := http.HTTPDelete(Config.APIEndpoint+"/api/v1/order/"+orderID,"?ts="+strconv.FormatInt(time.Now().Unix()*1000,10))
	if err != nil {
		return responses.CancelResponse{}, err
	}

	err = json.Unmarshal(res, &errorResponse)
	if err != nil {
		return responses.CancelResponse{}, err
	}

	if len(errorResponse.Error) > 0 {
		return responses.CancelResponse{}, errors.New(errorResponse.Error)
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
func Account() (responses.Balances, error){
	var balances responses.Balances
	var errorResponse responses.ErrorResponse


	http.Config = Config
	res, err := http.HTTPGet(Config.APIEndpoint+"/api/v1/account","?ts="+strconv.FormatInt(time.Now().Unix()*1000,10),true)
	if err != nil {
		return responses.Balances{}, err
	}

	err = json.Unmarshal(res, &errorResponse)
	if err != nil {
		return responses.Balances{}, err
	}

	if len(errorResponse.Error) > 0 {
		return responses.Balances{}, errors.New(errorResponse.Error)
	}


	err = json.Unmarshal(res, &balances)
	if err != nil {
		return responses.Balances{}, err
	}

	return balances, nil
}

