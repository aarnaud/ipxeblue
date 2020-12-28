package midlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"strings"
)

func MidlewareDevWebUI() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Next()
			return
		}
		if c.Request.URL.Path == "/ui" || c.Request.URL.Path == "/ui/" || c.Request.URL.Path == "/sockjs-node" {
			ProxyDevWebUI(c)
		}
		if strings.Contains(c.Request.Header.Get("Referer"), "/ui") {
			ProxyDevWebUI(c)
			return
		}
		if c.Request.Header.Get("Upgrade") == "websocket" {
			ProxyDevWebUI(c)
			return
		}
		c.Next()
	}
}

func ProxyDevWebUI(c *gin.Context) {
	director := func(req *http.Request) {
		req.URL = c.Request.URL
		req.URL.Scheme = "http"
		req.URL.Host = "localhost:3000"
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)
}
