package gocryptobot_tests

import (
	"github.com/florianpautot/go-arbitrage-trade-io/model/responses"
	"github.com/florianpautot/go-arbitrage-trade-io/utils"
	"github.com/golang/glog"
	"testing"
)

//TestGenerateSingature :
func TestGenerateSignature(t *testing.T) {
	config, err := utils.LoadConfig("../config.yaml")
	if err != nil {
		t.Error("Error while parsing config file:", err)
	}
	signature := utils.GenerateSignature(config.APIKey, config.APISecret)
	if signature != "cef45552612a60be0a1d58abfee61389cca3479acd00d697e2ccc4230a12bb8d999385032347280e9b6a8dd6005ab24f1b41ea7970aaff20066723e10d3b048b" {
		t.Error("Error, signature should be equal :", "cef45552612a60be0a1d58abfee61389cca3479acd00d697e2ccc4230a12bb8d999385032347280e9b6a8dd6005ab24f1b41ea7970aaff20066723e10d3b048b", " but got", signature)
	}
}

func TestLoadConfig(t *testing.T) {

	config, err := utils.LoadConfig("../config.yaml")
	if err != nil {
		t.Error("Error while parsing config file:", err)
	}

	if config.MinProfit != 1.003 {
		t.Error("Expected MinProfit to equal", 1, "but instead got ", config.MinProfit)
	}
	glog.Infoln(config)
}

func TestFormatBalances(t *testing.T) {
	balances := responses.Balances{
		Code:      0,
		Timestamp: 0,
		Balances: []responses.Balance{
			responses.Balance{
				Asset:     "btc",
				Available: "",
				Locked:    "",
			},
			responses.Balance{
				Asset:     "usdt",
				Available: "",
				Locked:    "",
			},
		},
	}

	formattedBalances := utils.FormatBalance(balances.Balances)
	if len(formattedBalances) != 2 {
		t.Error("Error, expected formattedBalances length of 2 but got", len(formattedBalances))
	}

	if formattedBalances["btc"] == (responses.Balance{}) {
		t.Error("Error, expected eth balance but got", formattedBalances["eth"])
	}

	if formattedBalances["usdt"] == (responses.Balance{}) {
		t.Error("Error, expected usdt balance but got", formattedBalances["eth"])
	}

	if formattedBalances["eth"] != (responses.Balance{}) {
		t.Error("Error, expected empty balance but got", formattedBalances["eth"])
	}
}

func TestFormatInfo(t *testing.T) {
	infos := responses.Infos{
		Code:      0,
		Timestamp: 0,
		Symbols: []responses.Symbol{
			responses.Symbol{
				Symbol:              "eth",
				Status:              "",
				BaseAsset:           "",
				BaseAssetPrecision:  0,
				QuoteAsset:          "",
				QuoteAssetPrecision: 0,
			},
			responses.Symbol{
				Symbol:              "btc",
				Status:              "",
				BaseAsset:           "",
				BaseAssetPrecision:  0,
				QuoteAsset:          "",
				QuoteAssetPrecision: 0,
			},
		},
	}

	formattedInfos := utils.FormatInfos(infos.Symbols)
	if len(formattedInfos) != 2 {
		t.Error("Error, expected formattedInfos length of 2 but got", len(formattedInfos))
	}

	if formattedInfos["eth"] == (responses.Symbol{}) {
		t.Error("Error, expected eth symbol info but got", formattedInfos["eth"])
	}

	if formattedInfos["btc"] == (responses.Symbol{}) {
		t.Error("Error, expected usdt symbol info but got", formattedInfos["eth"])
	}

	if formattedInfos["usdt"] != (responses.Symbol{}) {
		t.Error("Error, expected empty symbol info but got", formattedInfos["eth"])
	}
}

func TestFormatTickers(t *testing.T) {
	tickers := responses.Tickers{
		Tickers: []responses.Ticker{
			responses.Ticker{
				Symbol:      "btc_usdt",
				AskPrice:    "0",
				AskQty:      "0",
				BidPrice:    "0",
				BidQty:      "0",
				LastPrice:   "0",
				LastQty:     "0",
				Volume:      "0",
				QuoteVolume: "0",
				OpenTime:    0,
				CloseTime:   0,
			}, responses.Ticker{
				Symbol:      "eth_usdt",
				AskPrice:    "0",
				AskQty:      "0",
				BidPrice:    "0",
				BidQty:      "0",
				LastPrice:   "0",
				LastQty:     "0",
				Volume:      "0",
				QuoteVolume: "0",
				OpenTime:    0,
				CloseTime:   0,
			},
		},
	}

	formattedTickers, symbols := utils.FormatTickers(tickers.Tickers)
	if len(formattedTickers) != 2 {
		t.Error("Error, expected formattedTickers length of 2 but got", len(formattedTickers))
	}

	if len(symbols) != 2 {
		t.Error("Error, expected symbols length of 2 but got", len(symbols))
	}

	if formattedTickers["btc_usdt"] == (responses.Ticker{}) {
		t.Error("Error, expected btc ticker but got", formattedTickers["btc"])
	}

	if formattedTickers["eth_usdt"] == (responses.Ticker{}) {
		t.Error("Error, expected eth ticker but got", formattedTickers["btc"])
	}

	if formattedTickers["eth_btc"] != (responses.Ticker{}) {
		t.Error("Error, expected eth ticker but got", formattedTickers["btc"])
	}
}
