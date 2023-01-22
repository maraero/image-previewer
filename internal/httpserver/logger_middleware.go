package httpserver

import (
	"fmt"
	"net/http"

	"github.com/maraero/image-previewer/internal/logger"
)

func loggerMiddleware(next http.Handler, l logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		uri := r.URL.String()
		method := r.Method

		l.Info(fmt.Sprintf("%s %s %s", uri, method, userAgent(r)))
	})
}

func userAgent(r *http.Request) string {
	userAgents := r.Header["User-Agent"]
	if len(userAgents) > 0 {
		return "\"" + userAgents[0] + "\""
	}
	return ""
}
