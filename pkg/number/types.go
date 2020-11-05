package number

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	percentageRegex = regexp.MustCompile(`([\-0-9.]+)%`)
)

func String2Float(a string) float64 {
	f, _ := strconv.ParseFloat(a, 64)
	return f
}

// "0.96%" -> 0.96
func String2Percentage(a string) float64 {
	groups := percentageRegex.FindStringSubmatch(a)
	percentageStr := groups[1]
	return String2Float(percentageStr)
}

func Float2String(f float64, precise int) string {
	return fmt.Sprintf(fmt.Sprintf("%%.%df", precise), f)
}

func Float2StringWithSign(f float64, precise int) string {
	if f > 0 {
		return fmt.Sprintf("+%s", Float2String(f, precise))
	} else {
		return Float2String(f, precise)
	}
}
