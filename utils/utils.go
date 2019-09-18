package utils

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"github.com/adshao/go-binance"
	"github.com/florianpautot/go-arbitrage/global"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
)

//LoadConfig :
func LoadConfig(path string) (global.Config, error) {

	config := global.Config{}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return global.Config{}, err
	}
	err = yaml.UnmarshalStrict(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return global.Config{}, err

	}

	return config, nil
}

//GenerateSignature :
func GenerateSignature(input string, secret string) string {
	hmac512 := hmac.New(sha512.New, []byte(secret))
	hmac512.Write([]byte(input))
	signature := hex.EncodeToString(hmac512.Sum(nil))
	return signature;
}

func FormatBalance(balances []binance.Balance) map[string]binance.Balance {
	formattedBalance := make(map[string]binance.Balance)
	for _, balance := range balances {
		formattedBalance[balance.Asset] = balance
	}
	return formattedBalance
}

func FormatInfos(infos []binance.Symbol) map[string]binance.Symbol {
	formattedInfos := make(map[string]binance.Symbol)
	for i:=0; i < len(infos); i++ {
		formattedInfos[infos[i].Symbol] = infos[i]
	}
	return formattedInfos
}

func FormatTickers(tickers []*binance.BookTicker) (map[string]binance.BookTicker, []string) {
	formattedTickers := make(map[string]binance.BookTicker)
	existingAssets := make(map[string]bool)

	var symbols []string

	for _, ticker := range tickers {
		asset := ticker.Symbol[0:3]
		formattedTickers[ticker.Symbol] = *ticker
		if !strings.Contains(global.GlobalConfig.Exclusions, asset) && existingAssets[asset] == false {
			symbols = append(symbols, asset)
			existingAssets[asset] = true
		}
	}

	return formattedTickers, symbols
}

func RoundUp(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return
}

func RoundDown(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Floor(digit)
	newVal = round / pow
	return
}

func UpdateCachedBalances() {
	balances, err := global.Binance.NewGetAccountService().Do(context.Background())
	if err != nil {
		glog.V(1).Info(err.Error())
	}

	if len(balances.Balances) > 0 {
		formattedBalances := FormatBalance(balances.Balances)
		global.GlobalConfig.MaxBTC,err = strconv.ParseFloat(formattedBalances["btc"].Free,64)
		if err != nil {
			glog.V(1).Info(err.Error())
		}
		global.GlobalConfig.MaxUSDT,err = strconv.ParseFloat(formattedBalances["usdt"].Free,64)
		if err != nil {
			glog.V(1).Info(err.Error())
		}
		global.GlobalConfig.MaxETH,err = strconv.ParseFloat(formattedBalances["eth"].Free,64)
		if err != nil {
			glog.V(1).Info(err.Error())
		}
	}
}