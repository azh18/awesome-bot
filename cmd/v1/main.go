package main

import (
	"flag"
	"fmt"

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
	defer browser.GetSeleniumFactory().Stop()
	o := initOptions()
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

func initOptions() *options.Options {
	o := &options.Options{}
	flag.StringVar(&o.LarkConfigPath, "lark-config", "", "lark config if use")
	flag.Parse()
	fmt.Printf("options=%#v", o)
	return o
}
