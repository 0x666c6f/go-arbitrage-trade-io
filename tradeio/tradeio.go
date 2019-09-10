package tradeio

import (
	"encoding/json"
	"errors"
	"github.com/florianpautot/go-arbitrage-trade-io/http"
	"github.com/florianpautot/go-arbitrage-trade-io/model"
	"github.com/florianpautot/go-arbitrage-trade-io/model/requests"
	"github.com/florianpautot/go-arbitrage-trade-io/model/responses"
	errors2 "github.com/florianpautot/go-arbitrage-trade-io/model/responses/errors"
	"github.com/florianpautot/go-arbitrage-trade-io/utils"
	"github.com/golang/glog"
	"strconv"
	"time"
)


//Info :
func  Info() (responses.Infos, error){
	var infos responses.Infos
	var errorResponse errors2.ErrorResponse

	res, err := http.HTTPGet(model.GlobalConfig.APIEndpoint+"/api/v1/info","",false)
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
	var errorResponse errors2.ErrorResponse

	res, err := http.HTTPGet(model.GlobalConfig.APIEndpoint+"/api/v1/tickers","",false)
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
	var errorResponse errors2.OrderErrorResponse
	marshOrder, err := json.Marshal(order)
	res, err := http.HTTPPost(model.GlobalConfig.APIEndpoint+"/api/v1/order",marshOrder)
	if err != nil {
		return responses.OrderResponse{}, err
	}

	err = json.Unmarshal(res, &errorResponse)
	if err != nil {
		return responses.OrderResponse{}, err
	}

	if len(errorResponse.Errors) > 0 {
		return responses.OrderResponse{}, errors.New(errorResponse.Errors[0].Message)
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
	var errorResponse errors2.ErrorResponse

	res, err := http.HTTPDelete(model.GlobalConfig.APIEndpoint+"/api/v1/order/"+orderID,"?ts="+strconv.FormatInt(time.Now().Unix()*1000,10))
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
	var errorResponse errors2.ErrorResponse


	res, err := http.HTTPGet(model.GlobalConfig.APIEndpoint+"/api/v1/account","?ts="+strconv.FormatInt(time.Now().Unix()*1000,10),true)
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

func UpdateCachedBalances() {
	balances, err := Account()
	if err != nil {
		glog.V(1).Info(err.Error())
	}

	if len(balances.Balances) > 0 {
		formattedBalances := utils.FormatBalance(balances.Balances)
		model.GlobalConfig.MaxBTC,err = strconv.ParseFloat(formattedBalances["btc"].Available,64)
		if err != nil {
			glog.V(1).Info(err.Error())
		}
		model.GlobalConfig.MaxUSDT,err = strconv.ParseFloat(formattedBalances["usdt"].Available,64)
		if err != nil {
			glog.V(1).Info(err.Error())
		}
		model.GlobalConfig.MaxETH,err = strconv.ParseFloat(formattedBalances["eth"].Available,64)
		if err != nil {
			glog.V(1).Info(err.Error())
		}
	}
}
