package number

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString2Percentage(t *testing.T) {
	percentage := String2Percentage("-0.98%")
	assert.EqualValues(t, -0.98, percentage)
}
