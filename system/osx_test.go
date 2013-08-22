package system

import (
	"testing"
)

func TestHoge(t *testing.T) {
	stat := GetSystemStat()
	if stat != nil {
		return
	}
}
