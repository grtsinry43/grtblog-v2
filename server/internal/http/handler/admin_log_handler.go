package handler

import (
	"bufio"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/grtsinry43/grtblog-v2/server/internal/http/response"
)

// AdminLogHandler 读取应用日志（仅后台查看）。
type AdminLogHandler struct {
	logPath  string
	maxLines int
}

func NewAdminLogHandler(logPath string, maxLines int) *AdminLogHandler {
	if maxLines <= 0 {
		maxLines = 200
	}
	return &AdminLogHandler{logPath: logPath, maxLines: maxLines}
}

type LogLinesEnvelope struct {
	Code   int           `json:"code"`
	BizErr string        `json:"bizErr"`
	Msg    string        `json:"msg"`
	Data   []string      `json:"data"`
	Meta   response.Meta `json:"meta"`
}

// List godoc
// @Summary 查看最近日志
// @Tags Admin-Log
// @Produce json
// @Success 200 {object} LogLinesEnvelope
// @Security BearerAuth
// @Router /admin/logs [get]
func (h *AdminLogHandler) List(c *fiber.Ctx) error {
	file, err := os.Open(h.logPath)
	if err != nil {
		return response.NewBizErrorWithMsg(response.ServerError, "无法读取日志文件")
	}
	defer file.Close()

	lines := tailLines(file, h.maxLines)
	return response.Success(c, lines)
}

// tailLines 读取文件末尾最多 max 行（简单扫描）。
func tailLines(file *os.File, max int) []string {
	scanner := bufio.NewScanner(file)
	buf := make([]string, 0, max)
	for scanner.Scan() {
		buf = append(buf, scanner.Text())
		if len(buf) > max {
			buf = buf[1:]
		}
	}
	return buf
}
