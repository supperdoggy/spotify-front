package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/supperdoggy/webproject/frontend/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type IHandlers interface {
	InitHandlers()
}

type Handlers struct {
	logger *zap.Logger
	r *gin.Engine
	s service.IService
}

func NewHandlers(l *zap.Logger, r *gin.Engine, s *service.IService) IHandlers {
	return &Handlers{logger: l, r: r, s: s}
}

func (h *Handlers) InitHandlers() {
	h.r.GET("/", h.index)

	apiv1 := h.r.Group("/api/v1")
	{
		apiv1.GET("/getallsongs", h.getAllSongs)
	}
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
