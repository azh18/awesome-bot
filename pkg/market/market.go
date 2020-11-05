package market

import (
	"fmt"
	"sort"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zbw0046/awesome-bot/pkg/message"
	"k8s.io/klog/v2"
)

var (
	browserKey = "xueqiu"
)

type OverviewModule struct {
	messageSender message.Sender
}

func NewOverviewModule(messageSender message.Sender) *OverviewModule {
	return &OverviewModule{messageSender: messageSender}
}

func (m *OverviewModule) Start(stop <-chan struct{}) {
	m.do()
	crontab := cron.New()
	crontab.AddFunc("CRON_TZ=Asia/Shanghai 35 9 * * *", func() {
		m.do()
	})
	crontab.AddFunc("CRON_TZ=Asia/Shanghai 0 10 * * *", func() {
		m.do()
	})
	crontab.AddFunc("CRON_TZ=Asia/Shanghai 0 12 * * *", func() {
		m.do()
	})
	crontab.AddFunc("CRON_TZ=Asia/Shanghai 0 14 * * *", func() {
		m.do()
	})
	crontab.AddFunc("CRON_TZ=Asia/Shanghai 0 15 * * *", func() {
		m.do()
	})
	crontab.Start()
}

func (m *OverviewModule) Run(message.Sender, string) error {
	return m.do()
}

func (m *OverviewModule) do() error {
	klog.Infof("trigger flush overview module")
	title := fmt.Sprintf("当前股市概况（%s）", time.Now().Format("2006-01-02 15:04:05"))
	msgObj := &message.Message{
		Title: title,
	}

	var codeLists []string
	codeLists = append(codeLists, MainIndexCodeList...)
	codeLists = append(codeLists, ETFCodeList...)
	codeLists = append(codeLists, LOFCodeList...)

	stocksInfo, err := getStockInformation(codeLists)
	if err != nil {
		klog.Errorf("get stock information error: %s", err.Error())
		return err
	}

	// 大盘
	mainBoardBlock := &message.Block{
		Title: "大盘概况",
	}
	for _, code := range MainIndexCodeList {
		mainBoardBlock.Lines = append(mainBoardBlock.Lines,
			fmt.Sprintf("%s(%s): %s", CodeChineseMap[code], code, stocksInfo.Data[code]))
	}

	// ETF
	etfBlock := &message.Block{
		Title: "ETF概况",
	}
	sort.Slice(ETFCodeList, func(i, j int) bool {
		return stocksInfo.Data[ETFCodeList[i]].RaisePct > stocksInfo.Data[ETFCodeList[j]].RaisePct
	})
	for _, code := range ETFCodeList {
		etfBlock.Lines = append(etfBlock.Lines,
			fmt.Sprintf("%s(%s): %s", CodeChineseMap[code], code, stocksInfo.Data[code]))
	}

	// LOF
	lofBlock := &message.Block{
		Title: "LOF概况",
	}
	for _, code := range LOFCodeList {
		lofBlock.Lines = append(lofBlock.Lines,
			fmt.Sprintf("%s(%s): %s", CodeChineseMap[code], code, stocksInfo.Data[code]))
	}

	msgObj.Content = append(msgObj.Content, mainBoardBlock, etfBlock, lofBlock)
	err = m.messageSender.SendMessage(msgObj)
	if err != nil {
		klog.Errorf("send message error: %s", err.Error())
	}
	return err
}
