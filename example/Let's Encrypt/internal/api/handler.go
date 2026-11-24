package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"lets/internal/model"
	"lets/internal/service"
)

// Handler API 处理器
type Handler struct {
	svc *service.CertService
}

func NewHandler(svc *service.CertService) *Handler {
	return &Handler{svc: svc}
}

// Apply 申请证书
func (h *Handler) Apply(c *gin.Context) {
	var req model.ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}
	resp, err := h.svc.Apply(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}
	ok(c, resp)
}

// Renew 手动续期
func (h *Handler) Renew(c *gin.Context) {
	var req struct {
		Domain string `json:"domain" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}
	if err := h.svc.Renew(c.Request.Context(), req.Domain); err != nil {
		handleError(c, err)
		return
	}
	ok(c, gin.H{"message": "renewal completed", "domain": req.Domain})
}

// List 证书列表
func (h *Handler) List(c *gin.Context) {
	resp, err := h.svc.List(c.Request.Context())
	if err != nil {
		internalError(c, err)
		return
	}
	ok(c, resp)
}

// Get 证书详情
func (h *Handler) Get(c *gin.Context) {
	resp, err := h.svc.Get(c.Request.Context(), c.Param("domain"))
	if err != nil {
		handleError(c, err)
		return
	}
	ok(c, resp)
}

// Delete 删除证书
func (h *Handler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("domain")); err != nil {
		handleError(c, err)
		return
	}
	ok(c, gin.H{"message": "deleted"})
}

// Logs 操作日志
func (h *Handler) Logs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	logs, err := h.svc.ListLogs(limit)
	if err != nil {
		internalError(c, err)
		return
	}
	ok(c, logs)
}

// DNSProviders 列出已配置的 DNS 服务商
func (h *Handler) DNSProviders(c *gin.Context) {
	ok(c, gin.H{"providers": h.svc.ListDNSProviders()})
}

func ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data})
}

func badRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": msg})
}

func internalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrCertExists):
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": err.Error()})
	case errors.Is(err, service.ErrCertNotFound):
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
	default:
		internalError(c, err)
	}
}
