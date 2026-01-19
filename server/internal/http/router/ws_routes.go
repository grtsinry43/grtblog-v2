package router

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/http/handler"
	"github.com/grtsinry43/grtblog-v2/server/internal/ws"
)

func registerWSRoutes(v2 fiber.Router, manager *ws.Manager) {
	wsHandler := handler.NewWSHandler(manager)

	v2.Use("/ws", func(c *fiber.Ctx) error {
		if !websocket.IsWebSocketUpgrade(c) {
			return fiber.ErrUpgradeRequired
		}

		roomKey, err := parseWSRoomKey(c)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		c.Locals("wsRoomKey", roomKey)
		return c.Next()
	})

	v2.Get("/ws", websocket.New(wsHandler.Handle))
}

func parseWSRoomKey(c *fiber.Ctx) (string, error) {
	roomType := strings.TrimSpace(c.Query("type"))
	if roomType == "" {
		return "", fmt.Errorf("missing room type")
	}
	switch roomType {
	case "article", "moment", "page":
	default:
		return "", fmt.Errorf("invalid room type")
	}

	idStr := strings.TrimSpace(c.Query("id"))
	if idStr == "" {
		return "", fmt.Errorf("missing room id")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return "", fmt.Errorf("invalid room id")
	}

	return fmt.Sprintf("%s:%d", roomType, id), nil
}
