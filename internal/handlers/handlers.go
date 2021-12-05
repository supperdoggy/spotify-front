package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/spotify-web-project/spotify-front/internal/service"
	"github.com/supperdoggy/spotify-web-project/spotify-front/internal/structs"
	"go.uber.org/zap"
	"net/http"
)

type IHandlers interface {
	InitHandlers()
}

type Handlers struct {
	logger *zap.Logger
	r      *gin.Engine
	s      service.IService
}

func NewHandlers(l *zap.Logger, r *gin.Engine, s *service.IService) IHandlers {
	return &Handlers{logger: l, r: r, s: *s}
}

func (h *Handlers) InitHandlers() {
	h.r.GET("/", h.index)
	h.r.GET("/upload", h.downloadFile)
	h.r.GET("/auth", h.authPage)

	apiv1 := h.r.Group("/api/v1")
	{
		apiv1.GET("/getallsongs", h.getAllSongs)
		apiv1.POST("/song", h.uploadSong)
		apiv1.POST("/auth", h.auth)
	}
}

func (h *Handlers) authPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
	return
}

func (h *Handlers) index(c *gin.Context) {
	c.HTML(http.StatusOK, "play.html", gin.H{})
	return
}

func (h *Handlers) getAllSongs(c *gin.Context) {
	respose, err := h.s.GetAllSongs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, respose)
}

func (h *Handlers) downloadFile(c *gin.Context) {
	c.HTML(http.StatusOK, "download.html", gin.H{})
}

func (h *Handlers) uploadSong(c *gin.Context) {
	var data structs.UploadSongRequest
	if err := c.Bind(&data); err != nil {
		h.logger.Error("error binding req", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant read your file man"})
		return
	}

	// todo forward data to backend
	err := h.s.UploadNewSong(&data)
	if err != nil {
		h.logger.Info("error when creating new song", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handlers) auth(c *gin.Context) {
	var req structs.AuthReq
	var resp structs.AuthResp
	if err := c.Bind(&req); err != nil {
		h.logger.Error("error binding request", zap.Error(err))
		resp.Error = "error binding request"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp, err := h.s.Auth(req)
	if err != nil {
		h.logger.Error("error when Auth()", zap.Error(err))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}
