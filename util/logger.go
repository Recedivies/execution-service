package util

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// HTTPLogger logs a gin HTTP request in JSON format. Uses the
// default logger from rs/zerolog.
func HttpLogger() gin.HandlerFunc {
	return StructuredLogger(&log.Logger)
}

type ResponseRecorder struct {
	gin.ResponseWriter
	Body []byte
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func StructuredLogger(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Start timer
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		param := gin.LogFormatterParams{}
		rec := &ResponseRecorder{ResponseWriter: c.Writer}

		c.Writer = rec

		// Process request
		c.Next()

		// Fill the params
		param.Latency = time.Since(start) // Stop timer
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		// Log using the params
		var logEvent *zerolog.Event
		if param.StatusCode >= 400 {
			logEvent = logger.Error().Bytes("body", rec.Body)
		} else {
			logEvent = logger.Info()
		}

		logEvent.Str("client_id", param.ClientIP).
			Str("method", param.Method).
			Int("status_code", param.StatusCode).
			Str("status_text", http.StatusText(param.StatusCode)).
			Str("path", param.Path).
			Str("latency", param.Latency.String()).
			Msg(param.ErrorMessage)
	}
}
