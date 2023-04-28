package service

import (
	"context"

	pb "fmnx.io/dev/repo/cmd/generated/proto/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func UnaryLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if req, ok := req.(*pb.UploadRequest); ok {
			req.Content = nil
		}
		if err != nil {
			logrus.Errorf("Request [%s] error, req: [%s], resp: [%s], err: [%s]", info.FullMethod, req, resp, err.Error())
			return resp, err
		}
		logrus.Infof("Request [%s] success, req: [%s], resp: [%s]", info.FullMethod, req, resp)
		return resp, err
	}
}
