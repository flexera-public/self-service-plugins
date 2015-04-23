package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// Simple HTTP logger middleware
func HttpLogger(logger *log.Logger) echo.Middleware {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			msg := fmt.Sprintf(`Processing GET "%s"`, c.Request.URL.String())
			originIp := c.Request.Header.Get("X-Forwarded-For")
			if originIp == "" {
				originIp = c.Request.Header.Get("X-Originating-IP")
			}
			if originIp != "" {
				msg += fmt.Sprintf(" (for %s)", originIp)
			}
			if reqId := c.Get("RequestID"); reqId != nil {
				msg += fmt.Sprintf(" - Request ID: %v", reqId)
			}
			logger.Print(msg)
			start := time.Now()
			err := h(c)
			if err != nil {
				return err
			}
			elapsed := time.Since(start)
			var status string
			var size int
			if resp := c.Response; resp != nil {
				status = http.StatusText(resp.Status())
				size = resp.Size()
			}
			logger.Printf(`Completed in %s | %s | %d bytes`, elapsed, status, size)
			return nil
		}
	}
}
