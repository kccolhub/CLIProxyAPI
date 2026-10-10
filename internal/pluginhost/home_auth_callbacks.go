package pluginhost

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/router-for-me/CLIProxyAPI/v8/internal/home"
)

func (h *Host) homeAuthCallback(ctx context.Context, method, index string) ([]byte, bool, error) {
	h.mu.Lock()
	enabled := h.runtimeConfig != nil && h.runtimeConfig.Home.Enabled
	h.mu.Unlock()
	if !enabled {
		return nil, false, nil
	}
	client := home.Current()
	if client == nil {
		return nil, true, fmt.Errorf("Home credential service unavailable")
	}
	raw, errRead := client.ReadPluginAuth(ctx, h.callbackCallerPluginID(ctx, ""), method, index)
	if errRead != nil {
		return nil, true, errRead
	}
	envelope, errMarshal := marshalRPCResult(json.RawMessage(raw))
	return envelope, true, errMarshal
}
