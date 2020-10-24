package arbitrage

import (
	"testing"

	"github.com/zbw0046/awesome-bot/pkg/message"
	"github.com/zbw0046/awesome-bot/pkg/test"
)

func TestArbitrage(t *testing.T) {
	test.InitAll()
	sender := message.GetSender(message.SenderLark)
	module := Module{}
	module.Run(sender, "")
}
