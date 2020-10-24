package starter

import "github.com/zbw0046/awesome-bot/pkg/message"

type Module interface {
	Start(stop <-chan struct{})
	Run(sender message.Sender, extra string) error
}
