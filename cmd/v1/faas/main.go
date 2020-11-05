package main

import (
	"encoding/json"
	"fmt"
	"time"

	gr "github.com/awesome-fc/golang-runtime"
	"github.com/zbw0046/awesome-bot/cmd/v1/options"
	"github.com/zbw0046/awesome-bot/pkg/arbitrage"
	"github.com/zbw0046/awesome-bot/pkg/browser"
	"github.com/zbw0046/awesome-bot/pkg/market"
	"github.com/zbw0046/awesome-bot/pkg/message"
	"github.com/zbw0046/awesome-bot/pkg/message/lark"
	"github.com/zbw0046/awesome-bot/pkg/starter"
	"github.com/zbw0046/awesome-bot/pkg/utils/graceful"
)

var (
	moduleMap map[string]starter.Module
)

/*
{
    "triggerTime":"2018-02-09T05:49:00Z",
    "triggerName":"my_trigger",
    "payload":"awesome-fc"
}
*/

type TriggerRequest struct {
	TriggerTime time.Time `json:"triggerTime"`
	TriggerName string    `json:"triggerName"`
	Payload     string    `json:"payload"`
}

type Request struct {
	ModuleCallList []*moduleCall `json:"moduleCallList"`
}

type moduleCall struct {
	Module string `json:"module"`
	Sender string `json:"sender"` // such as: lark, dingtalk
	Extra  string `json:"extra"`  // extra params for different module
}

type Response struct {
	Result map[string]*result `json:"result"`
	ErrMsg string             `json:"err_msg"`
}

type result struct {
	ErrMsg   string `json:"err_msg"`
	Response string `json:"response"`
}

func handler(ctx *gr.FCContext, event []byte) ([]byte, error) {
	fcLogger := gr.GetLogger().WithField("requestId", ctx.RequestID)
	_, err := json.Marshal(ctx)
	if err != nil {
		fcLogger.Error("error:", err)
	}
	fcLogger.Infof("hello golang!")

	req := TriggerRequest{}
	resp := Response{
		Result: map[string]*result{},
	}
	if err := json.Unmarshal(event, &req); err != nil {
		resp.ErrMsg = err.Error()
		fcLogger.Errorf("unmarshal request error: %s", err.Error())
		return marshal(resp), nil
	}

	fcLogger.Infof("trigger info: %#v", req)

	realReq := &Request{}
	if err := json.Unmarshal([]byte(req.Payload), realReq); err != nil {
		resp.ErrMsg = err.Error()
		fcLogger.Errorf("unmarshal request error: %s", err.Error())
		return marshal(resp), nil
	}

	for _, moduleCall := range realReq.ModuleCallList {
		fcLogger.Infof("module called: %s", moduleCall.Module)
		subResult := runModule(moduleCall)
		resp.Result[moduleCall.Module] = subResult
	}
	return marshal(resp), nil
}

func runModule(moduleCall *moduleCall) (subResult *result) {
	subResult = &result{}
	if moduleCall.Sender == "" {
		moduleCall.Sender = message.SenderLark
	}
	sender := message.GetSender(moduleCall.Sender)
	if sender == nil {
		subResult.ErrMsg = fmt.Sprintf("cannot get sender: %s", moduleCall.Sender)
		return
	}

	err := moduleMap[moduleCall.Module].Run(sender, moduleCall.Extra)
	if err != nil {
		subResult.ErrMsg = err.Error()
		return
	}
	subResult.Response = "ok"
	return
}

func marshal(object interface{}) []byte {
	result, _ := json.Marshal(object)
	return result
}

func main() {
	stopCh := graceful.SetupSignalHandler()
	o := options.InitOptions()
	browser.InitSelenium(o.ChromeDriverPath, o.ChromeDriverPort)
	defer browser.GetSeleniumFactory().Stop()
	// init lark config
	lark.InitApp(o.LarkConfigPath)
	moduleMap = map[string]starter.Module{
		"overview":  market.NewOverviewModule(&message.LarkSender{}),
		"arbitrage": &arbitrage.Module{},
	}
	gr.Start(handler, nil)
	<-stopCh
}
