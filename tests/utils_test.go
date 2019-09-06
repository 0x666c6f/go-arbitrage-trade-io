package gocryptobot_tests

import (
	"fmt"
	"github.com/florianpautot/utils"
	"testing"
)

//TestGenerateSingature :
func TestGenerateSignature(t *testing.T) {
	config, err:= utils.LoadConfig("../config.yaml")
	if err != nil {
		t.Error("Error while parsing config file:",err)
	}
	signature := utils.GenerateSignature(config.APIKey,config.APISecret)
	if signature != "cef45552612a60be0a1d58abfee61389cca3479acd00d697e2ccc4230a12bb8d999385032347280e9b6a8dd6005ab24f1b41ea7970aaff20066723e10d3b048b"{
		t.Error("Error, signature should be equal :","cef45552612a60be0a1d58abfee61389cca3479acd00d697e2ccc4230a12bb8d999385032347280e9b6a8dd6005ab24f1b41ea7970aaff20066723e10d3b048b"," but got",signature)
	}
}

func TestLoadConfig(t *testing.T){

	config, err:= utils.LoadConfig("../config.yaml")
	if err != nil {
		t.Error("Error while parsing config file:",err)
	}

	if config.MinProfit != 1.003 {
		t.Error("Expected MinProfit to equal",1,"but instead got ", config.MinProfit)
	}
	fmt.Println(config)
}
