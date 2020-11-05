package test

import (
	"github.com/zbw0046/awesome-bot/cmd/v1/options"
	"github.com/zbw0046/awesome-bot/pkg/browser"
	"github.com/zbw0046/awesome-bot/pkg/message/lark"
)

func InitAll() {
	o := &options.Options{
		LarkConfigPath:   "/Users/bytedance/code/awesome-bot/conf/lark_config.yaml",
		ChromeDriverPath: "/Users/bytedance/code/awesome-bot/chromedriver",
		ChromeDriverPort: 9090,
	}
	browser.InitSelenium(o.ChromeDriverPath, o.ChromeDriverPort)
	defer browser.GetSeleniumFactory().Stop()
	// init lark config
	lark.InitApp(o.LarkConfigPath)
}
