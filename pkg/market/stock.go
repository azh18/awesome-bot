package market

import (
	"fmt"
	"regexp"

	"github.com/tebeka/selenium"
	"github.com/zbw0046/awesome-bot/pkg/browser"
	"github.com/zbw0046/awesome-bot/pkg/number"
	"k8s.io/klog"
)

const (
	xueQiuStockUrl = "https://xueqiu.com/S/%s"
)

var (
	regexStockPrice, _  = regexp.Compile(`([0-9\-+.]+)`)
	regexStockChange, _ = regexp.Compile(`([0-9\-+.]+)\s+([0-9\-+.]+)`)
)

type stocksInfo struct {
	Data                   map[string]*stockInfo
	ForeignMoneyIn         float64
	NationalDebtRevenue10Y float64
}

type stockInfo struct {
	precise  int
	Current  float64
	Raise    float64
	RaisePct float64

	PE float64
	PB float64
}

func (s *stockInfo) String() string {
	return fmt.Sprintf("%s  %s(%s%%)  PE: %.2f  PB:%.2f",
		number.Float2String(s.Current, s.precise),
		number.Float2StringWithSign(s.Raise, s.precise),
		number.Float2StringWithSign(s.RaisePct, 2), s.PE, s.PB)
}

func getStockInformation(codes []string) (*stocksInfo, error) {
	stocksInfo := &stocksInfo{
		Data:                   map[string]*stockInfo{},
		ForeignMoneyIn:         0,
		NationalDebtRevenue10Y: 0,
	}
	for _, code := range codes {
		info, err := getStockInformationFromXueQiu(code)
		if err != nil {
			return nil, fmt.Errorf("get stockInfo of %s(%s) error: %s", code, CodeChineseMap[code], err.Error())
		}
		klog.Infof("get %s success", code)
		stocksInfo.Data[code] = info
	}
	return stocksInfo, nil
}

func getStockInformationFromXueQiu(code string) (*stockInfo, error) {
	info := &stockInfo{}
	wd := browser.GetSeleniumFactory().GetSelenium(browserKey)
	if wd == nil {
		return nil, fmt.Errorf("get selenium error")
	}

	url := fmt.Sprintf(xueQiuStockUrl, code)
	// todo: add retry
	err := wd.WebDriver.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch page %s error: %s", url, err.Error())
	}

	stockCurrentText, err := getTextValueFromPage(wd, `//div[@class='stock-current']`)
	if err != nil {
		return nil, fmt.Errorf("get stock current value in TEXT error: %s", err.Error())
	}
	groups := regexStockPrice.FindStringSubmatch(stockCurrentText)
	info.Current = number.String2Float(groups[1])

	stockRaise, err := getTextValueFromPage(wd, `//div[@class='stock-change']`)
	if err != nil {
		return nil, fmt.Errorf("get stock raise value in TEXT error: %s", err.Error())
	}

	groups = regexStockChange.FindStringSubmatch(stockRaise)
	if len(groups) != 3 {
		return nil, fmt.Errorf("run regex on value %s error: %v.need 3 fields", stockRaise, groups)
	}

	info.Raise = number.String2Float(groups[1])
	info.RaisePct = number.String2Float(groups[2])

	switch CodeType(code) {
	case CodeTypeETF, CodeTypeLOF:
		info.precise = 3
	default:
		info.precise = 2
	}

	return info, nil
}

func getTextValueFromPage(wd *browser.SeleniumService, xpath string) (string, error) {
	elem, err := wd.WebDriver.FindElement(selenium.ByXPATH, xpath)
	if err != nil {
		return "", fmt.Errorf("get value by XPATH error: %s", err.Error())
	}
	text, err := elem.Text()
	if err != nil {
		return "", fmt.Errorf("get value in TEXT error: %s", err.Error())
	}
	return text, nil
}
