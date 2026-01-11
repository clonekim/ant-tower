package server

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:embed dist
var staticFiles embed.FS

func SetupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	skipPaths := []string{"/api", "/ws", "/assets"}

	// 헬퍼 함수: 경로가 스킵 대상인지 확인
	shouldSkipSpa := func(path string) bool {
		for _, prefix := range skipPaths {
			if strings.HasPrefix(path, prefix) {
				return true
			}
		}
		return false
	}

	api := router.Group("/api")
	{
		api.GET("/uptime", h.GetUptime)
		api.GET("/login", h.GetCurrentUser)
		api.GET("/process", h.GetProcessList)
		api.POST("/process/kill", h.KillProcess)
		api.POST("/power", h.ControlPower)
	}

	router.GET("/ws", h.Hub.HandleWebSocket)

	distFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to access embedded 'dist'")
	}

	assetsFS, err := fs.Sub(distFS, "assets")
	if err == nil {
		router.StaticFS("/assets", http.FS(assetsFS))
	} else {
		log.Warn().Msg("Assets folder not found")
	}

	indexData, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		log.Error().Err(err).Msg("index.html missing")
		indexData = []byte("<h1>Frontend Build Not Found</h1>")
	}

	serveIndex := func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexData)
	}

	router.GET("/", serveIndex)

	router.NoRoute(func(c *gin.Context) {
		if shouldSkipSpa(c.Request.URL.Path) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		serveIndex(c)
	})

	return router
}
