package lark

import (
	"context"
	"fmt"
	"testing"

	"github.com/larksuite/botframework-go/SDK/message"
	"github.com/larksuite/botframework-go/SDK/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/zbw0046/awesome-bot/pkg/message/lark/builder"
	"k8s.io/klog/v2"
)

func TestApp(t *testing.T) {
	msg := builder.NewBuilder("test msg")
	msg.AddThing(&builder.Thing{
		Type:  builder.TextThingType,
		Lines: []string{"testline1, testline2"},
	})
	card, err := msg.Build()
	assert.Nil(t, err)
	InitApp("/Users/bytedance/code/awesome-bot/conf/lark_config.yaml")
	resp, err := message.SendTextMessage(context.Background(), "", larkConfig.AppID,
		&protocol.UserInfo{
			ID:   "zhangbowen.alf@bytedance.com",
			Type: protocol.UserTypeEmail,
		}, "", "tests")
	fmt.Printf("resp=%#v\n", resp)
	assert.Nil(t, err)

	resp, cardResp, err := SendTo(context.Background(), "me", card)
	klog.Infof("resp=%#v, cardResp=%#v, err=%#v", resp, cardResp, err)
	assert.Nil(t, err)
	klog.Flush()
}
