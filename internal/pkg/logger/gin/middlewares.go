package gin

import (
	"github.com/artemys/pprof-visualizer/internal/pkg/logger"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Ginzap Code from https://github.com/gin-contrib/zap/blob/master/zap.go :
// + zap.logger -> utils.Logger
// + excludedUrls to skip some urls
func Ginzap(logger *logger.Logger, timeFormat string, utc bool, alwaysLog bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			if alwaysLog || c.Writer.Status() >= 300 {
				logger.Info(path,
					zap.Int("status", c.Writer.Status()),
					zap.String("method", c.Request.Method),
					zap.String("path", path),
					zap.String("query", query),
					zap.String("ip", c.ClientIP()),
					zap.String("user-agent", c.Request.UserAgent()),
					zap.String("time", end.Format(timeFormat)),
					zap.Duration("latency", latency),
				)
			}
		}
	}
}

func RecoveryWithZap(logger *logger.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
