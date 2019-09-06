package model

type Config struct {
	APIKey string `yaml:"APIKey"`
	APISecret string `yaml:"APISecret"`
	APIEndpoint string `yaml:"APIEndpoint"`
	MinBTC float32 `yaml:"MinBTC"`
	MaxBTC float32 `yaml:"MaxBTC"`
	MinETH float32 `yaml:"MinETH"`
	MaxETH float32 `yaml:"MaxETH"`
	MinUSDT float32 `yaml:"MinUSDT"`
	MaxUSDT float32 `yaml:"MaxUSDT"`
	MinProfit float32 `yaml:"MinProfit"`
	OrderMinuteLimit int32 `yaml:"OrderMinuteLimit"`
	APIMinuteLimit int32 `yaml:"APIMinuteLimit"`
	Timeout int32 `yaml:"Timeout"`
	Fees float32 `yaml:"Fees"`
	StartSecond int32 `yaml:"StartSecond"`
	EndSecond int32 `yaml:"EndSecond"`
	Debug bool `yaml:"Debug"`
	Exclusions []string `yaml:"Exclusions"`
}
