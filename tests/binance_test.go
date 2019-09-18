package gocryptobot_tests

import (
	"github.com/florianpautot/go-arbitrage/global"
	"github.com/florianpautot/go-arbitrage/utils"

	"testing"
)


func TestBinance(t *testing.T){
	config, err:= utils.LoadConfig("../config.yaml")
	if err != nil {
		return
	}
	global.GlobalConfig = config


}