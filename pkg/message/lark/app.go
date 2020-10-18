package lark

import (
	"context"
	"fmt"

	"github.com/larksuite/botframework-go/SDK/appconfig"
	"github.com/larksuite/botframework-go/SDK/message"
	"github.com/larksuite/botframework-go/SDK/protocol"
	"k8s.io/klog/v2"
)

const (
	tenantKey = ""
)

var (
	chatMap    = map[string]*protocol.UserInfo{}
	larkConfig *appconfig.AppConfig
)

func InitApp(confPath string) {
	config := ParseConfig(confPath)

	botConfig := config.BotConfig
	larkConfig = &appconfig.AppConfig{
		AppID:       botConfig.AppID,       // get it from lark-voucher and basic information。
		AppType:     protocol.InternalApp,  // AppType only has two types: Independent Software Vendor App（ISVApp） or Internal App.
		AppSecret:   botConfig.AppSecret,   // get it from lark-voucher and basic information.
		VerifyToken: botConfig.VerifyToken, // get it from lark-event subscriptions.
		EncryptKey:  "",                    // get it from lark-event subscriptions.
	}
	appconfig.Init(*larkConfig)
	c, _ := appconfig.GetConfig(botConfig.AppID)

	for _, userItem := range config.Users {
		fmt.Printf("user=%#v\n", userItem)
		chatMap[userItem.Name] = &protocol.UserInfo{
			ID:   userItem.Email,
			Type: protocol.UserTypeEmail,
		}
		// chat.CreateChat(context.Background(), tenantKey, botConfig.AppID, &protocol.CreateChatRequest{
		// 	Name:          "财经机器人",
		// 	Description:   "财经机器人",
		// 	UserIDs:       []string{},
		// 	OpenIDs:       nil,
		// 	ChatI18nNames: nil,
		// })
	}

	klog.Infof("lark config=%#v", c)
	// _ = event.EventRegister(config.GlobalConf.BotConfig.AppID, protocol.EventTypeMessage, EventMessage)
	// _ = event.BotRecvMsgRegister(constant.AppID, "help", BotRecvMsgHelp)
	// _ = event.BotRecvMsgRegister(constant.AppID, "rs", MessageWrapper(svc.GetResource))
	// _ = event.BotRecvMsgRegister(constant.AppID, "dj", MessageWrapper(svc.GetTaskManager))
	// _ = event.BotRecvMsgRegister(constant.AppID, "f", MessageWrapper(svc.SearchTask))
	// _ = event.BotRecvMsgRegister(constant.AppID, "ld", MessageWrapper(svc.GetLeader))
	// _ = event.BotRecvMsgRegister(constant.AppID, "status", MessageWrapper(svc.GetStatus))
	// _ = event.BotRecvMsgRegister(constant.AppID, "summarize", MessageWrapper(svc.SummarizeTask))
	// event.IgnoreSign(constant.AppID, true)
	// _ = event.CardRegister(constant.AppID, "seeSomeTask", CardCallbackWrapper(svc.ActionSeeSomeTasks))
	// _ = event.CardRegister(constant.AppID, "maintenance", CardCallbackWrapper(svc.ActionSeeSomeTasks))
	// _ = event.CardRegister(constant.AppID, "laijigecuowurizhi", TenMoreHandler)
}

type eventUser struct {
	User string `json:"user"`
	Text string `json:"text"`
}

func SendTo(ctx context.Context, name string, ret interface{}) (*protocol.SendMsgResponse, *protocol.SendCardMsgResponse, error) {
	if ch, ok := chatMap[name]; !ok {
		return nil, nil, fmt.Errorf("unknown chat name: %s in: %+v", name, chatMap)
	} else {
		return sendMessage(ctx, ch, "", ret)
	}
}

func sendMessage(ctx context.Context, user *protocol.UserInfo, openMessageID string, ret interface{}) (r *protocol.SendMsgResponse, r2 *protocol.SendCardMsgResponse, err error) {
	switch ret.(type) {
	case string:
		{
			if openMessageID != "" {
				r, err = message.SendTextMessage(ctx, tenantKey, larkConfig.AppID, user, openMessageID, prefixAt(ctx, ret.(string)))
			} else {
				r, err = message.SendTextMessage(ctx, tenantKey, larkConfig.AppID, user, openMessageID, ret.(string))
			}
		}
	case map[protocol.Language]*protocol.RichTextForm:
		{
			r, err = message.SendRichTextMessage(ctx, tenantKey, larkConfig.AppID, user, openMessageID, ret.(map[protocol.Language]*protocol.RichTextForm))
		}
	case *protocol.CardForm:
		{
			r2, err = message.SendCardMessage(ctx, tenantKey, larkConfig.AppID, user, openMessageID, *ret.(*protocol.CardForm), true)
		}
	}

	return r, r2, err
}

func prefixAt(ctx context.Context, text string) string {
	userKey := "ff"
	userID, ok := ctx.Value(userKey).(string)
	klog.Infof("userID: %s ok: %v", userID, ok)
	return fmt.Sprintf(`<at user_id="%s"></at>
%s`, userID, text)
}
