package gocryptobot

import (
	"fmt"
	"github.com/florianpautot/go-arbitrage-trade-io/arbitrage"
	"github.com/florianpautot/go-arbitrage-trade-io/http"
	"github.com/florianpautot/go-arbitrage-trade-io/tradeio"
	"github.com/florianpautot/go-arbitrage-trade-io/utils"
)

func main() {
	config,err := utils.LoadConfig("config.yaml")
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	http.Config = config
	tradeio.Config = config
	arbitrage.Config = config

	infos,err := tradeio.Info()
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	arbitrage.Infos = utils.FormatInfos(infos.Symbols)

	arbitrage.Start()
}
