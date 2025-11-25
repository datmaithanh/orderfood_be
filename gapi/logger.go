package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLoger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()
	results, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "gRPC").
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Str("method", info.FullMethod).
		Msg("gRPC request received")
	return results, err
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body	   []byte
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	rw.body = body
	return rw.ResponseWriter.Write(body)
}

func HTTPLoger(handle http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rec := &responseWriter{ResponseWriter: res, statusCode: http.StatusOK}
		handle.ServeHTTP(rec, req)
		duration := time.Since(startTime)

		logger := log.Info()
		if rec.statusCode != http.StatusOK {
			logger = log.Error().Bytes("response_body", rec.body)
		}
		logger.Str("protocol", "HTTP").
		Str("method", req.Method).
		Str("path", req.RequestURI).
		Int("status_code", rec.statusCode).
		Str("status_text", http.StatusText(rec.statusCode)).
		Dur("duration", duration).
		Msg("HTTP request received")
	})
}
