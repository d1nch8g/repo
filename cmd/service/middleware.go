package service

import (
	"context"
	"fmt"
	"strings"

	pb "fmnx.io/dev/repo/cmd/generated/proto/v1"

	"google.golang.org/grpc"
)

func UnaryLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if req, ok := req.(*pb.UploadRequest); ok {
			req.Content = nil
		}
		if err != nil {
			fmt.Printf("Request [%s] error, req: [%s], resp: [%s], err: [%s]\n", info.FullMethod, req, resp, err.Error())
			return resp, err
		}
		fmt.Printf("Request [%s] success, req: [%s], resp: [%s]\n", info.FullMethod, req, resp)
		return resp, err
	}
}

func Validate(in string) bool {
	return !strings.ContainsAny(in, "<>'\"#$&\\^*?")
}
