package rpc

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/nash567/GoSentinel/api/v1/pb/goSentinel"
	application "github.com/nash567/GoSentinel/internal/service/application/model"
	"github.com/nash567/GoSentinel/internal/service/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	goSentinel.UnimplementedGoSentinelServiceServer
	applicationSvc application.Service
	authSvc        auth.Service
	Key            string
}

func NewServer(applicationSvc application.Service, authSvc auth.Service, key string) *Server {
	return &Server{
		applicationSvc: applicationSvc,
		authSvc:        authSvc,
		Key:            key,
	}
}

func (s *Server) SendVerifcationNotification(ctx context.Context, req *goSentinel.SendApplicationNotificationRequest) (*goSentinel.SendApplicationNotificationResponse, error) {
	token, err := s.applicationSvc.SendVerifcationNotification(ctx, req.Email, req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("failed to send  notification: %v", err))
	}
	return &goSentinel.SendApplicationNotificationResponse{
		Token: aws.StringValue(token),
	}, nil
}

func (s *Server) VerifyApplication(ctx context.Context, req *goSentinel.VerifyApplicationRequest) (*emptypb.Empty, error) {
	err := s.applicationSvc.VerifyApplication(ctx, req.Key)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("failed to verify application: %v", err))
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) GetApplicationSecrets(ctx context.Context, req *emptypb.Empty) (*goSentinel.GetApplicationSecretResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("failed to verify application"))
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("failed to verify application"))
	}
	res, err := s.applicationSvc.GetApplicationSecret(ctx, authHeader[0])
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("failed to verify application: %v", err))
	}
	secret, err := s.authSvc.DecryptData(res.Secret, s.Key)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to decrypt secret: %v", err))
	}
	return &goSentinel.GetApplicationSecretResponse{
		ClientID:     res.ID,
		ClientSecret: string(secret),
	}, nil
}
