package plugin

import "testing"

func TestRegister(t *testing.T) {
	var plugin Plugin
	Register("test", plugin)
}
