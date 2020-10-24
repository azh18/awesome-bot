package arbitrage

import "github.com/zbw0046/awesome-bot/pkg/message"

/*
Run at 7:30 on each transaction day.
Compute premium rate of 华宝油气(162411)
*/
type Module struct {
}

func (m *Module) Start(stop <-chan struct{}) {
	panic("implement me")
}

func (m *Module) Run(sender message.Sender, extra string) error {
	Analyse162411(sender)
	return nil
}
