package market

import (
	"testing"

	"github.com/zbw0046/awesome-bot/pkg/message"
	"github.com/zbw0046/awesome-bot/pkg/test"
)

func TestMarketModule(t *testing.T) {
	test.InitAll()
	sender := message.GetSender(message.SenderLark)
	module := NewOverviewModule(sender)
	module.Run(sender, "afternoon")
}
