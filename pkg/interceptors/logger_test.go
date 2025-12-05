package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
)

func TestLoggerInterceptorAddsLoggerToContext(t *testing.T) {
	l := logger.NewSlogLogger()
	interceptor := LoggerInterceptor(l)
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/Call"}

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		require.Equal(t, l, logger.FromContext(ctx))
		return "ok", nil
	}

	resp, err := interceptor(context.Background(), nil, info, handler)
	require.NoError(t, err)
	require.Equal(t, "ok", resp)
}



