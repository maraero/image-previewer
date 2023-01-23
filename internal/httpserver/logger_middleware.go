package httpserver

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/maraero/image-previewer/internal/logger"
)

func loggerMiddleware(next http.Handler, l logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rww := newResponseWriterWrapper(w)
		defer func() {
			l.Info(fmt.Sprintf(
				"[Remote addr: %s] [URL: %s] [Method: %s] [Status code : %d] [User-agent: %s] [Execution time: %v]",
				requestAddr(r),
				r.URL.String(),
				r.Method,
				*(rww.statusCode),
				userAgent(r),
				time.Since(start)),
			)
		}()
		next.ServeHTTP(rww, r)
	})
}

func requestAddr(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func userAgent(r *http.Request) string {
	userAgents := r.Header["User-Agent"]
	if len(userAgents) > 0 {
		return "\"" + userAgents[0] + "\""
	}
	return ""
}
