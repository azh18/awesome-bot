package arbitrage

import (
	"fmt"
	"time"

	"github.com/zbw0046/awesome-bot/pkg/browser"
	"github.com/zbw0046/awesome-bot/pkg/message"
	"github.com/zbw0046/awesome-bot/pkg/number"
	"github.com/zbw0046/awesome-bot/pkg/utils"
	"k8s.io/klog"
)

const (
	NFYY = "501018"
	HBYQ = "162411"
)

const (
	browserKey = "arbitrage"
)

const (
	url162411 = "https://palmmicro.com/woody/res/sz162411cn.php"
)

func Analyse162411(sender message.Sender) {
	basicBlock := message.Block{
		Title: "162411基本信息",
	}
	wd := browser.GetSeleniumFactory().GetSelenium(browserKey).WebDriver
	if err := wd.Get(url162411); err != nil {
		klog.Errorf("fetch 162411 data from %s error: %s", url162411, err.Error())
		return
	}
	estimatedValueStr, err := browser.ExtractTextValue(wd, `//*[@id="estimation"]/tbody/tr[2]/td[3]/font`)
	if err != nil {
		klog.Errorf("get estimated value error: %s", err.Error())
		return
	}
	lofValueStr, err := browser.ExtractTextValue(wd, `//*[@id="reference"]/tbody/tr[2]/td[2]/font`)
	if err != nil {
		klog.Errorf("get LOF value error: %s", err.Error())
		return
	}
	lastValueStr, err := browser.ExtractTextValue(wd, `//*[@id="estimation"]/tbody/tr[2]/td[2]`)
	if err != nil {
		klog.Errorf("get last day value error: %s", err.Error())
		return
	}
	currentTime := time.Now()

	estimatedValue, lastDayValue, lofValue := number.String2Float(estimatedValueStr), number.String2Float(lastValueStr), number.String2Float(lofValueStr)
	estimatedPremium := (lofValue - estimatedValue) / estimatedValue * 100

	basicBlock.Lines = []string{
		fmt.Sprintf("时间：%v", utils.FormatTime(currentTime)),
		fmt.Sprintf("昨日净值：%.3f, 今晨收盘估值：%.3f，当前盘中交易价：%.3f，溢价率：%.2f%%",
			lastDayValue, estimatedValue, lofValue, estimatedPremium),
		"\n",
	}

	// todo: analyse whether there is an arbitrage chance
	arbitrarySuggestBlock := message.Block{
		Title: "交易建议",
	}
	haveArbitrageChance := false
	// T-1溢价率和T-2溢价率
	t1PremiumStr, err := browser.ExtractTextValue(wd, `//*[@id="SZ162411fundhistory"]/tbody/tr[2]/td[4]/font`)
	if err != nil {
		klog.Errorf("Error: %s", err.Error())
		return
	}
	t2PremiumStr, err := browser.ExtractTextValue(wd, `//*[@id="SZ162411fundhistory"]/tbody/tr[3]/td[4]/font`)
	if err != nil {
		klog.Errorf("Error: %s", err.Error())
		return
	}
	t1Premium, t2Premium := number.String2Percentage(t1PremiumStr), number.String2Percentage(t2PremiumStr)
	if t1Premium > 3 && t2Premium > 3 {
		haveArbitrageChance = true
		arbitrarySuggestBlock.Lines = append(arbitrarySuggestBlock.Lines,
			fmt.Sprintf("建议套利：T-1溢价率：%.2f%%，T-2溢价率：%.2f%%，均超过3%%", t1Premium, t2Premium))
	}

	// 前一天XOP是否大跌
	xopRaiseStr, err := browser.ExtractTextValue(wd, `//*[@id="reference"]/tbody/tr[3]/td[3]/font`)
	if err != nil {
		klog.Errorf("Error: %s", err.Error())
		return
	}
	xopRaise := number.String2Percentage(xopRaiseStr)
	if xopRaise < -3 {
		haveArbitrageChance = true
		arbitrarySuggestBlock.Lines = append(arbitrarySuggestBlock.Lines,
			fmt.Sprintf("建议套利：上一日XOP跌幅为%.2f%%，超过3%%。观察今日各项指标，如果今夜不会大涨，可以套利", xopRaise))
	}

	// send a message
	bannerBlock := message.Block{
		Title: "额外提醒",
		Lines: []string{
			"交易费提醒：加和约0.2%",
			"可以观察油价，若交易时段相比当日8点出现大跌（2%）以上，可以赶快场内卖出然后申购。（TODO：自动化监控，faas可能不够）",
		},
	}

	// 估值说明：https://palmmicro.com/woody/blog/entertainment/20150818cn.php
	msg := &message.Message{
		Title: "162411提醒",
		Content: []*message.Block{
			&basicBlock, &arbitrarySuggestBlock, &bannerBlock,
		},
		Links: map[string]string{
			"162411数据来源": url162411,
		},
	}

	// 如果有套利机会，消息加急
	if haveArbitrageChance {
		msg.IsUrgent = true
	}

	if err := sender.SendMessage(msg); err != nil {
		klog.Errorf("send message error: %s", err.Error())
	}
}
