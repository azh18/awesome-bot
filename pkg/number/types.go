package number

import (
	"fmt"
	"strconv"
)

func String2Float(a string) float64 {
	f, _ := strconv.ParseFloat(a, 64)
	return f
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
