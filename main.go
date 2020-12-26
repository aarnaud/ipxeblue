package main

import (
	"fmt"
	"github.com/aarnaud/ipxeblue/config"
	"github.com/aarnaud/ipxeblue/controllers"
	_ "github.com/aarnaud/ipxeblue/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// @title ipxeblue API
// @version 0.1
// @description Manage PXE boot
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv() // To get the value from the config file using key// viper package read .env

	port := 8080
	if p := viper.GetInt("PORT"); p != 0 {
		port = p
	}
	router := gin.Default()

	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"x-total-count", "content-length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.LoadHTMLGlob("templates/*")
	db := config.Database()

	// Provide db variable to controllers
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	if gin.Mode() == gin.DebugMode {
		// Configure SwaggerUI
		url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
		// proxies UI call to nodejs react server
		router.Use(MidlewareWebUI())
	} else {
		// Serve react-admin webui
		router.Static("/ui", "./ui")
	}

	// iPXE request
	router.GET("/", controllers.Index)

	// API
	v1 := router.Group("/api/v1")
	{
		v1.GET("/computers", controllers.ListComputers)
		v1.GET("/computers/:id", controllers.GetComputer)
		v1.PUT("/computers/:id", controllers.UpdateComputer)
		v1.DELETE("/computers/:id", controllers.DeleteComputer)
	}

	router.Run(fmt.Sprintf(":%d", port))
}

func MidlewareWebUI() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Next()
			return
		}
		if c.Request.URL.Path == "/ui" || c.Request.URL.Path == "/ui/" {
			ProxyWebUI(c)
		}
		if strings.Contains(c.Request.Header.Get("Referer"), "/ui") {
			ProxyWebUI(c)
			return
		}
		if c.Request.Header.Get("Upgrade") == "websocket" {
			ProxyWebUI(c)
			return
		}
		c.Next()
	}
}

func ProxyWebUI(c *gin.Context) {
	director := func(req *http.Request) {
		req.URL = c.Request.URL
		req.URL.Scheme = "http"
		req.URL.Host = "localhost:3000"
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)
}
