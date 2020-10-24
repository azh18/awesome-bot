package message

import (
	"fmt"

	"github.com/larksuite/botframework-go/SDK/protocol"
	"github.com/zbw0046/awesome-bot/pkg/message/lark/builder"
)

type Message struct {
	Title    string
	Content  []*Block
	Links    map[string]string // text -> url
	IsUrgent bool
}

func (m *Message) ToLarkCard() (*protocol.CardForm, error) {
	msg := builder.NewBuilder(m.Title)
	for _, block := range m.Content {
		lines := []string{fmt.Sprintf("**%s**", block.Title), "\n"}
		lines = append(lines, block.Lines...)
		msg.AddThing(&builder.Thing{
			Type:  builder.TextThingType,
			Lines: lines,
		}, builder.HrThing)
	}

	msg.AddThing(&builder.Thing{
		Type:  builder.TextThingType,
		Lines: []string{"[雪球行情]($xueqiu)", "[华宝油气系列数据]($huabao)"},
		URLMap: map[string]string{
			"xueqiu": "https://xueqiu.com/S/SH000001",
			"huabao": "https://palmmicro.com/woody/res/sz162411cn.php",
		},
	})
	return msg.Build()
}
