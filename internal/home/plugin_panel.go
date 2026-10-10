package home

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"fmt"
)

func (c *Client) ReadPluginAuth(ctx context.Context, pluginID, method, index string) ([]byte, error) {
	if pluginID == "" {
		return nil, fmt.Errorf("plugin caller identity unavailable")
	}
	cmd, errClient := c.commandClient()
	if errClient != nil {
		return nil, errClient
	}
	payload, _ := json.Marshal(map[string]string{"plugin_id": pluginID, "method": method, "auth_index": index})
	raw, errRead := cmd.Do(ctx, "GET", "plugin-auth", string(payload)).Text()
	return []byte(raw), errRead
}

func (c *Client) RegisterPluginPanel(ctx context.Context, port int) error {
	c.mu.Lock()
	if c.pluginPanelToken == "" {
		c.pluginPanelToken = rand.Text() + rand.Text()
	}
	token := c.pluginPanelToken
	c.mu.Unlock()
	cmd, errClient := c.commandClient()
	if errClient != nil {
		return errClient
	}
	payload, _ := json.Marshal(map[string]any{"port": port, "token": token})
	return cmd.Do(ctx, "GET", "plugin-panel", string(payload)).Err()
}

func (c *Client) VerifyPluginPanelToken(token string) bool {
	if c == nil || c.dispatchFenced.Load() {
		return false
	}
	c.mu.Lock()
	expected := c.pluginPanelToken
	c.mu.Unlock()
	return expected != "" && subtle.ConstantTimeCompare([]byte(expected), []byte(token)) == 1
}

type pluginCallerKey struct{}

// WithPluginCaller marks a host-originated execution; HTTP headers cannot set it.
func WithPluginCaller(ctx context.Context, id string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, pluginCallerKey{}, id)
}

func pluginCaller(ctx context.Context) string {
	id, _ := ctx.Value(pluginCallerKey{}).(string)
	return id
}
