package main

import (
	"fmt"
	"github.com/aarnaud/ipxeblue/controllers"
	_ "github.com/aarnaud/ipxeblue/docs"
	"github.com/aarnaud/ipxeblue/midlewares"
	"github.com/aarnaud/ipxeblue/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/pin/tftp/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
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
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	appconf := utils.GetConfig()

	router := gin.New()
	router.Use(logger.SetLogger(logger.Config{
		SkipPath: []string{"/healthz"},
	}))
	router.Use(gin.Recovery())

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
	db := utils.Database()
	filestore := utils.NewFileStore(appconf)
	log.Info().Msg("starting TokenCleaner in goroutine")
	go utils.TokenCleaner(db)

	// Provide db variable to controllers
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("filestore", filestore)
		c.Set("config", appconf)
		c.Header("Cache-Control", "no-store, max-age=0")
		c.Next()
	})

	if gin.Mode() == gin.DebugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		// Configure SwaggerUI
		url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
		// proxies UI call to nodejs react server
		router.Use(midlewares.MidlewareDevWebUI())
	} else {
		// Serve react-admin webui
		router.Static("/admin", "./admin")
	}

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	if appconf.GrubSupportEnabled {
		// Grub request without auth
		log.Info().Msg("enabling grub http endpoint")
		router.GET("/grub/", controllers.GrubScript)
	}

	// grub don't support long http queries
	// using TFTP to pass positional metadata in tftp path
	if appconf.TFTPEnabled {
		s := tftp.NewServer(controllers.GetTFTPReader(appconf, db), controllers.GetTFTPWriter(appconf))
		s.SetTimeout(30 * time.Second)
		go func() {
			log.Info().Msg("starting tftp server")
			err := s.ListenAndServe(":69")
			if err != nil {
				fmt.Fprintf(os.Stdout, "server: %v\n", err)
			}
		}()
	}

	// iPXE request with auth
	ipxeroute := router.Group("/", midlewares.BasicAuthIpxeAccount(false))
	ipxeroute.GET("/", controllers.IpxeScript)

	router.HEAD("/files/public/:uuid/*filepath", controllers.DownloadPublicFile)
	router.GET("/files/public/:uuid/*filepath", controllers.DownloadPublicFile)
	router.HEAD("/files/token/:token/:uuid/*filepath", controllers.DownloadProtectedFile)
	router.GET("/files/token/:token/:uuid/*filepath", controllers.DownloadProtectedFile)

	var v1 *gin.RouterGroup
	if appconf.EnableAPIAuth {
		// API
		v1 = router.Group("/api/v1", midlewares.BasicAuthIpxeAccount(true))
	} else {
		// API
		v1 = router.Group("/api/v1")
	}

	// Computer
	v1.GET("/computers", controllers.ListComputers)
	v1.GET("/computers/:id", controllers.GetComputer)
	v1.PUT("/computers/:id", controllers.UpdateComputer)
	v1.DELETE("/computers/:id", controllers.DeleteComputer)

	// IPXE account
	v1.GET("/ipxeaccounts", controllers.ListIpxeaccount)
	v1.GET("/ipxeaccounts/:username", controllers.GetIpxeaccount)
	v1.POST("/ipxeaccounts", controllers.CreateIpxeaccount)
	v1.PUT("/ipxeaccounts/:username", controllers.UpdateIpxeaccount)
	v1.DELETE("/ipxeaccounts/:username", controllers.DeleteIpxeaccount)

	// Bootentry
	v1.GET("/bootentries", controllers.ListBootentries)
	v1.GET("/bootentries/:uuid", controllers.GetBootentry)
	v1.POST("/bootentries", controllers.CreateBootentry)
	v1.PUT("/bootentries/:uuid", controllers.UpdateBootentry)
	v1.DELETE("/bootentries/:uuid", controllers.DeleteBootentry)
	v1.POST("/bootentries/:uuid/files/:name", controllers.UploadBootentryFile)
	v1.GET("/bootentries/:uuid/files/:name", controllers.DownloadBootentryFile)

	log.Info().Msg("starting http server")
	router.Run(fmt.Sprintf(":%d", appconf.Port))
}
