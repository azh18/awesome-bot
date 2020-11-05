package message

import (
	"context"
	"fmt"

	"github.com/zbw0046/awesome-bot/pkg/message/lark"
)

type Sender interface {
	SendMessage(message *Message) error
}

func GetSender(senderType string) Sender {
	switch senderType {
	case SenderLark:
		return &LarkSender{}
	}
	return nil
}

type LarkSender struct {
}

func (l *LarkSender) SendMessage(message *Message) error {
	card, err := message.ToLarkCard()
	if err != nil {
		return fmt.Errorf("ToLarkCard error: %s", err.Error())
	}
	_, _, err = lark.SendTo(context.Background(), "me", card, message.IsUrgent)
	return err
}
