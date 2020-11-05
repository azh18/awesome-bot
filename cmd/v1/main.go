package main

import (
	"github.com/zbw0046/awesome-bot/cmd/v1/options"
	"github.com/zbw0046/awesome-bot/pkg/browser"
	"github.com/zbw0046/awesome-bot/pkg/market"
	"github.com/zbw0046/awesome-bot/pkg/message"
	"github.com/zbw0046/awesome-bot/pkg/message/lark"
	"github.com/zbw0046/awesome-bot/pkg/starter"
	"github.com/zbw0046/awesome-bot/pkg/utils/graceful"
)

func main() {
	stopCh := graceful.SetupSignalHandler()
	o := options.InitOptions()
	browser.InitSelenium(o.ChromeDriverPath, o.ChromeDriverPort)
	defer browser.GetSeleniumFactory().Stop()
	// init lark config
	lark.InitApp(o.LarkConfigPath)
	modules := []starter.Module{
		market.NewOverviewModule(&message.LarkSender{}),
	}
	for _, module := range modules {
		module.Start(stopCh)
	}
	<-stopCh
}
