package server

import "testing"

func TestGetMainConfig(t *testing.T) {
	mainConfig := GetMainConfig("../test_data/")
	if mainConfig.Port == "" {
		t.Errorf("PORT is not define")
	}
	if mainConfig.UrlPath == "" {
		t.Errorf("PATH is not define")
	}
}
