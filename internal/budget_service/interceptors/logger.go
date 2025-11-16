
package interceptors

import (
    "context"
    "google.golang.org/grpc"
    "github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
)

func LoggerInterceptor(l logger.Logger) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (resp interface{}, err error) {
        ctx = logger.WithLogger(ctx, l)
        return handler(ctx, req)
    }
}
