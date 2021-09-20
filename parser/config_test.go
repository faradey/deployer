package parser

import "testing"

func TestGetConfig(t *testing.T) {
	lines := GetConfig("../test_data/")
	if len(lines) == 0 {
		t.Errorf("File deployer-config is empty")
	}
}
