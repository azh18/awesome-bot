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
	o := initOptions()
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

func initOptions() *options.Options {
	o := &options.Options{}
	flag.StringVar(&o.LarkConfigPath, "lark-config", "", "lark config if use")
	flag.StringVar(&o.ChromeDriverPath, "chrome-driver-path", "", "chrome driver path if used")
	flag.IntVar(&o.ChromeDriverPort, "chrome-driver-port", 4444, "chrome driver port if used (default 4444)")
	flag.Parse()
	fmt.Printf("options=%#v", o)
	return o
}
