package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/chirikova/go-anti-brute-force/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func loggingMiddleware(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		r interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		startTime := time.Now()
		response, err := handler(ctx, r)

		var ip, userAgent string
		var statusCode codes.Code
		peerInfo, ok := peer.FromContext(ctx)
		if ok {
			ip = peerInfo.Addr.String()
		}

		incomingContext, ok := metadata.FromIncomingContext(ctx)
		if ok {
			userAgent = incomingContext.Get("user-agent")[0]
		}

		if fromError, ok := status.FromError(err); ok {
			statusCode = fromError.Code()
		}

		logger.Info(
			fmt.Sprintf(
				"%s %s %v %d %s",
				ip,
				info.FullMethod,
				statusCode,
				time.Since(startTime),
				userAgent,
			),
		)

		return response, err
	}
}
