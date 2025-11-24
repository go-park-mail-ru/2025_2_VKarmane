package interceptors

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
)

func LoggerInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Info("gRPC request",
			"method", info.FullMethod,
		)

		resp, err := handler(ctx, req)

		if err != nil {
			logger.Error("gRPC request failed",
				"method", info.FullMethod,
				"error", err,
			)
		} else {
			logger.Info("gRPC request succeeded",
				"method", info.FullMethod,
			)
		}

		return resp, err
	}
}

