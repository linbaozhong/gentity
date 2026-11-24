package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter 创建 HTTP 路由
func NewRouter(h *Handler, webDir string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API 路由
	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/cert/apply", h.Apply)
		apiGroup.POST("/cert/renew", h.Renew)
		apiGroup.GET("/cert/list", h.List)
		apiGroup.GET("/cert/:domain", h.Get)
		apiGroup.DELETE("/cert/:domain", h.Delete)
		apiGroup.GET("/logs", h.Logs)
		apiGroup.GET("/dns/providers", h.DNSProviders)
	}

	// 管理界面（单文件静态页面）
	if webDir != "" {
		r.NoRoute(func(c *gin.Context) {
			c.File(webDir + "/index.html")
		})
		r.GET("/", func(c *gin.Context) {
			c.File(webDir + "/index.html")
		})
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
