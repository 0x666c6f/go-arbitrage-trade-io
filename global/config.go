package global

import binance "github.com/adshao/go-binance"

type Config struct {
	APIKey string `yaml:"APIKey"`
	APISecret string `yaml:"APISecret"`
	APIEndpoint string `yaml:"APIEndpoint"`
	MinBTC float64 `yaml:"MinBTC"`
	MaxBTC float64 `yaml:"MaxBTC"`
	MinETH float64 `yaml:"MinETH"`
	MaxETH float64 `yaml:"MaxETH"`
	MinUSDT float64 `yaml:"MinUSDT"`
	MaxUSDT float64 `yaml:"MaxUSDT"`
	MinProfit float64 `yaml:"MinProfit"`
	OrderMinuteLimit int `yaml:"OrderMinuteLimit"`
	APIMinuteLimit int `yaml:"APIMinuteLimit"`
	Timeout string `yaml:"Timeout"`
	Fees float64 `yaml:"Fees"`
	StartSecond int `yaml:"StartSecond"`
	EndSecond int `yaml:"EndSecond"`
	Debug bool `yaml:"Debug"`
	Exclusions string `yaml:"Exclusions"`
}

var GlobalConfig Config
var Binance *binance.Client
