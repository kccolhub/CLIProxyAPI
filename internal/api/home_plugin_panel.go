package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v8/internal/home"
)

func (s *Server) serveHomePluginPanel(c *gin.Context) {
	if s.cfg == nil || !s.cfg.Home.Enabled || s.pluginHost == nil {
		c.AbortWithStatus(404)
		return
	}
	client := home.Current()
	if !client.VerifyPluginPanelToken(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")) {
		c.AbortWithStatus(401)
		return
	}
	path := c.Param("path")
	if path == "/plugins" && c.Request.Method == http.MethodGet {
		plugins := make([]gin.H, 0)
		for _, info := range s.pluginHost.RegisteredPlugins() {
			menus := make([]gin.H, 0, len(info.Menus))
			for _, menu := range info.Menus {
				menus = append(menus, gin.H{"path": menu.Path, "menu": menu.Menu, "description": menu.Description})
			}
			plugins = append(plugins, gin.H{"id": info.ID, "name": info.Metadata.Name, "version": info.Metadata.Version, "menus": menus})
		}
		c.JSON(200, gin.H{"plugins": plugins})
		return
	}
	request := c.Request.Clone(c.Request.Context())
	request.URL.Path, request.URL.RawPath = path, ""
	request.Header.Del("Authorization")
	if strings.HasPrefix(path, "/v0/resource/plugins/") && (request.Method == http.MethodGet || request.Method == http.MethodHead) {
		if s.pluginHost.ServeResourceHTTP(c.Writer, request) {
			return
		}
	} else if strings.HasPrefix(path, "/v0/management/") {
		if s.pluginHost.ServeManagementHTTP(c.Writer, request) {
			return
		}
	}
	c.AbortWithStatus(404)
}
