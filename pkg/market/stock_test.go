package market

import (
	"fmt"
	"regexp"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestStock(t *testing.T) {
	re1, err := regexp.Compile(`([0-9\-+.]+)\s+([0-9\-+.]+)`)
	assert.Nil(t, err)
	fmt.Printf("%v", re1.FindStringSubmatch("+4.18 +0.13%"))

	info, err := getStockInformationFromXueQiu("SH000001")
	if err != nil {
		assert.Nil(t, err)
	}
	fmt.Printf("SH000001: %v", info)
}
