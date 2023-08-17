package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func CheckUUIDRequest(expectedUUID string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		uuid := getUUIDFromMetadata(ctx)

		if uuid == expectedUUID {
			resp, err = handler(ctx, req)
		}

		return resp, err
	}
}

func getUUIDFromMetadata(ctx context.Context) string {
	values := getValuesFromMetadata(ctx, "uuid")

	if len(values) < 1 {
		return ""
	}

	return values[0]
}

func getValuesFromMetadata(ctx context.Context, key string) []string {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return []string{}
	}

	return md.Get(key)
}
