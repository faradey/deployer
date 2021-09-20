package server

import "testing"

func TestGetMainConfig(t *testing.T) {
	mainConfig := GetMainConfig()
	if mainConfig.Dir != "" {
		t.Errorf("Dir is not define")
	}
}
