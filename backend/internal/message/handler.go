package message

import (
	"net/http"
	"strings"

	"feedsystem_video_go/internal/apierror"
	"feedsystem_video_go/internal/middleware/jwt"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) Send(c *gin.Context) {
	fromID, err := jwt.GetAccountID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var req SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.ToID == 0 || strings.TrimSpace(req.Content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to_id and content are required"})
		return
	}
	m := &Message{FromID: fromID, ToID: req.ToID, Content: req.Content}
	if err := h.service.Send(c.Request.Context(), m); err != nil {
		c.JSON(apierror.ClassifyHTTPStatus(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

func (h *Handler) List(c *gin.Context) {
	userID, err := jwt.GetAccountID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var req ListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.PeerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "peer_id is required"})
		return
	}
	msgs, err := h.service.List(c.Request.Context(), userID, req.PeerID, 50)
	if err != nil {
		c.JSON(apierror.ClassifyHTTPStatus(err), gin.H{"error": err.Error()})
		return
	}
	if msgs == nil {
		msgs = []Message{}
	}
	c.JSON(http.StatusOK, ListResponse{Messages: msgs})
}
