package options

import (
	"flag"
	"fmt"
)

type Options struct {
	LarkConfigPath   string
	ChromeDriverPath string
	ChromeDriverPort int
}

func InitOptions() *Options {
	o := &Options{}
	flag.StringVar(&o.LarkConfigPath, "lark-config", "", "lark config if use")
	flag.StringVar(&o.ChromeDriverPath, "chrome-driver-path", "", "chrome driver path if used")
	flag.IntVar(&o.ChromeDriverPort, "chrome-driver-port", 4444, "chrome driver port if used (default 4444)")
	flag.Parse()
	fmt.Printf("options=%#v", o)
	return o
}
