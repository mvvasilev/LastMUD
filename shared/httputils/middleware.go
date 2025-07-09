package httputils

import (
	"bytes"
	"code.haedhutner.dev/mvv/LastMUD/shared/log"
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"runtime/debug"
)

const LoggerContextKey = "httputils_logger"

type HTTPError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewHTTPError(status int, message string) *HTTPError {
	return &HTTPError{
		Status:  status,
		Message: message,
	}
}

func (e *HTTPError) Error() string {
	w := bytes.NewBuffer([]byte{})
	err := EncodeBody(w, e)

	if err != nil {
		return "" // TODO
	}

	return w.String()
}

func GetLogger(req *http.Request) *log.Logger {
	logger, ok := req.Context().Value(LoggerContextKey).(*log.Logger)

	if ok {
		return logger
	} else {
		return log.DefaultLogger
	}
}

func NewLoggingMiddleware(logger *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := context.WithValue(request.Context(), LoggerContextKey, logger)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

func NewErrorHandlingMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := GetLogger(r)

			defer func() {
				if err := recover(); err != nil {
					http.Error(w, NewHTTPError(http.StatusInternalServerError, "Internal Server Error").Error(), http.StatusInternalServerError)
					logger.Error("Unrecoverable error in request to ", r.RequestURI, ": ", err, "\n", string(debug.Stack()))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
