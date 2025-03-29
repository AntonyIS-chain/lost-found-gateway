package adapters

import (
	"github.com/AntonyIS-chain/lost-found-gateway/pkg"
	"github.com/gin-gonic/gin"
)

// RegisterProxyRoutes sets up reverse proxy routes for a given service
func RegisterProxyRoutes(router *gin.RouterGroup, target string) {
	router.Any("/*proxyPath", func(c *gin.Context) {
		proxyPath := c.Param("proxyPath")
		if proxyPath == "" {
			proxyPath = "/"
		}

		pkg.NewReverseProxy(target)(c)
	})
}
