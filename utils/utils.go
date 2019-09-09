package utils

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"github.com/florianpautot/go-arbitrage-trade-io/model"
	"github.com/florianpautot/go-arbitrage-trade-io/model/responses"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math"
	"strings"
)

//LoadConfig :
func LoadConfig(path string) (model.Config,error){

	config := model.Config{}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return model.Config{},err
	}
	err = yaml.UnmarshalStrict(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return model.Config{},err

	}

	return config,nil
}

//GenerateSignature :
func GenerateSignature(input string, secret string) string {
	hmac512 := hmac.New(sha512.New, []byte(secret))
	hmac512.Write([]byte(input))
	signature := hex.EncodeToString(hmac512.Sum(nil))
	return signature;
}

func FormatBalance(balances []responses.Balance) map[string]responses.Balance {
	formattedBalance := make(map[string]responses.Balance)
	for _,balance := range balances {
		formattedBalance[balance.Asset] = balance
	}
	return formattedBalance
}

func FormatInfos(infos []responses.Symbol) map[string]responses.Symbol {
	formattedInfos := make(map[string]responses.Symbol)
	for _,info := range infos {
		formattedInfos[info.Symbol] = info
	}
	return formattedInfos
}

func FormatTickers(tickers []responses.Ticker) (map[string]responses.Ticker, []string) {
	formattedTickers := make(map[string]responses.Ticker)
	existingTickers := make(map[string]bool)

	var symbols []string
	for _,ticker := range tickers {
		asset := strings.Split(ticker.Symbol,"_")[0]
		if !strings.Contains(model.GlobalConfig.Exclusions,asset) && existingTickers[asset] == false{
			formattedTickers[ticker.Symbol] = ticker
			symbols = append(symbols, asset)
			existingTickers[asset] = true
		}

	}
	return formattedTickers,symbols
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