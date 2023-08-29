package rpc

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/nash567/GoSentinel/api/v1/pb/goSentinel"
	application "github.com/nash567/GoSentinel/internal/service/application/model"
	"github.com/nash567/GoSentinel/internal/service/auth"
	authModel "github.com/nash567/GoSentinel/internal/service/auth/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	goSentinel.UnimplementedGoSentinelServiceServer
	applicationSvc application.Service
	authSvc        *auth.Service
	Key            string
}

func NewServer(applicationSvc application.Service, authSvc *auth.Service, key string) *Server {
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
func (s *Server) CreateApplicationSecret(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	err := s.applicationSvc.CreateApplicationIdentity(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("failed to create application identity: %v", err))
	}
	return &emptypb.Empty{}, nil
}
func (s *Server) GetApplicationSecret(ctx context.Context, _ *emptypb.Empty) (*goSentinel.Applicationcredentials, error) {

	res, err := s.applicationSvc.GetApplicationIdentity(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("failed to get application identity: %v", err))
	}

	return &goSentinel.Applicationcredentials{
		ApplicationID:     res.ApplicationID,
		ApplicationSecret: res.ApplicationSecret,
	}, nil
}

func (s *Server) GetApplicationToken(ctx context.Context, req *goSentinel.Applicationcredentials) (*goSentinel.GetApplicationTokenResponse, error) {
	token, err := s.authSvc.GetApplicationToken(ctx, authModel.Credentials{
		ApplicationID:     req.ApplicationID,
		ApplicationSecret: req.ApplicationSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("appliation not verified: %v", err)
	}
	return &goSentinel.GetApplicationTokenResponse{
		Token: aws.StringValue(token),
	}, nil
}
