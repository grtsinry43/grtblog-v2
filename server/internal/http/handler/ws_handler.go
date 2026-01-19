package handler

import (
	"github.com/gofiber/websocket/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/ws"
)

type WSHandler struct {
	manager *ws.Manager
}

func NewWSHandler(manager *ws.Manager) *WSHandler {
	return &WSHandler{manager: manager}
}

func (h *WSHandler) Handle(conn *websocket.Conn) {
	if h.manager == nil {
		return
	}
	roomKey, ok := conn.Locals("wsRoomKey").(string)
	if !ok || roomKey == "" {
		return
	}
	client, cached := h.manager.Join(roomKey, conn)
	if client == nil {
		return
	}
	for _, payload := range cached {
		_ = client.Write(payload)
	}

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	h.manager.Leave(roomKey, client)
}
