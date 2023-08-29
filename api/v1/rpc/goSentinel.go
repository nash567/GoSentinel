package rpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/nash567/GoSentinel/api/v1/pb/goSentinel"
	applicationModel "github.com/nash567/GoSentinel/internal/service/application/model"
	"github.com/nash567/GoSentinel/internal/service/auth"
	authModel "github.com/nash567/GoSentinel/internal/service/auth/model"
	userModel "github.com/nash567/GoSentinel/internal/service/user/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	goSentinel.UnimplementedGoSentinelServiceServer
	applicationSvc applicationModel.Service
	userSvc        userModel.Service
	authSvc        *auth.Service
	Key            string
}

func NewServer(applicationSvc applicationModel.Service, authSvc *auth.Service, encryptionKey string) *Server {
	return &Server{
		applicationSvc: applicationSvc,
		authSvc:        authSvc,
		Key:            encryptionKey,
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

func (s *Server) GetUser(ctx context.Context, req *goSentinel.GetUserRequest) (*goSentinel.GetUserResponse, error) {
	claims, ok := authModel.GetJWTClaimsFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	user, err := s.userSvc.GetUser(ctx, userModel.Filter{
		Email: []string{req.Email},
		ID:    []string{req.ID},
	}, claims.ApplicationID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return nil, status.Error(codes.NotFound, fmt.Sprintf("user do not exists please register: %v", err))
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("user not found: %v", err))
	}
	return &goSentinel.GetUserResponse{
		ID:    user.ID,
		Name:  user.Email,
		Email: user.Email,
	}, nil
}
func (s *Server) RegisterUser(ctx context.Context, req *goSentinel.RegisterUserRequest) (*emptypb.Empty, error) {
	claims, ok := authModel.GetJWTClaimsFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	err := s.userSvc.RegisterUser(ctx, userModel.User{
		Name:  req.Name,
		Email: req.Email,
	}, claims.ApplicationID)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("user cannot be created try after some time: %v", err))
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *goSentinel.UpdateUserRequest) (*emptypb.Empty, error) {
	claims, ok := authModel.GetJWTClaimsFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	err := s.userSvc.UpdateUser(ctx, userModel.UpdateUser{
		ID:       req.ID,
		Name:     req.Name,
		Password: req.Password,
	}, claims.ApplicationID)

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("user cannot be updated: %w", err))
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *goSentinel.DeleteUserRequest) (*emptypb.Empty, error) {
	claims, ok := authModel.GetJWTClaimsFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("claims not found in context")
	}
	err := s.userSvc.DeleteUser(ctx, req.ID, claims.ApplicationID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("user cannot be deleted: %v", err))
	}

	return &emptypb.Empty{}, nil

}
