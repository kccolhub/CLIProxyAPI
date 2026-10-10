package home

import (
	"context"
	"testing"
)

func TestPluginPanelTokenAndCallerIsolation(t *testing.T) {
	client := &Client{pluginPanelToken: "node-secret"}
	if !client.VerifyPluginPanelToken("node-secret") || client.VerifyPluginPanelToken("other-secret") {
		t.Fatal("invalid token verification")
	}
	client.dispatchFenced.Store(true)
	if client.VerifyPluginPanelToken("node-secret") {
		t.Fatal("fenced node accepted panel request")
	}
	if pluginCaller(context.Background()) != "" {
		t.Fatal("untrusted context has caller")
	}
	if pluginCaller(WithPluginCaller(context.Background(), "plugin")) != "plugin" {
		t.Fatal("host identity not propagated")
	}
}
