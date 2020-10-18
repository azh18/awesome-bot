package market

const (
	CodeSZZS = "SH000001"
	CodeSZCZ = "SZ399001"
	CodeCYB  = "SZ399006"
	CodeKC50 = "SH000688"

	CodeETF5G         = "SH515050"
	CodeETFQuanShang  = "SH512000"
	CodeETFXinNengChe = "SH515700"

	CodeLOFHBYouQi = "SZ162411"
)

var (
	MainIndexCodeList = []string{CodeSZZS, CodeSZCZ, CodeCYB, CodeKC50}
	ETFCodeList       = []string{CodeETF5G, CodeETFQuanShang, CodeETFXinNengChe}
	LOFCodeList       = []string{CodeLOFHBYouQi}
)

var (
	ETFIndexCodeMap = map[string]string{
		CodeETF5G:         "931079",
		CodeETFQuanShang:  "399975",
		CodeETFXinNengChe: "930997",
		CodeLOFHBYouQi:    "",
	}
)

var (
	CodeChineseMap = map[string]string{
		CodeSZZS:          "上证指数",
		CodeSZCZ:          "深证成指",
		CodeCYB:           "创业板",
		CodeKC50:          "科创50",
		CodeETF5G:         "5GETF",
		CodeETFQuanShang:  "券商ETF",
		CodeETFXinNengChe: "新能车",
		CodeLOFHBYouQi:    "华宝油气",
	}
)

const (
	CodeTypeMain = iota
	CodeTypeETF
	CodeTypeLOF
	CodeTypeIndex
	CodeTypeStock

	CodeTypeUnknown = -1
)

func CodeType(code string) int {
	switch code[2:4] {
	case "51":
		return CodeTypeETF
	case "00":
		if code == CodeSZZS {
			return CodeTypeMain
		} else {
			return CodeTypeStock
		}
	case "39":
		return CodeTypeMain
	case "16":
		return CodeTypeLOF
	}
	return CodeTypeUnknown
}
