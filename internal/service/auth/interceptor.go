package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nash567/GoSentinel/api/v1/pb/goSentinel"
	"github.com/nash567/GoSentinel/internal/service/auth/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Service) AuthenticationInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	grpcLog := grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
	grpclog.SetLoggerV2(grpcLog)

	start := time.Now()

	var (
		ctxWithClaims context.Context = ctx
		err           error
	)

	if info.FullMethod == "/goSentinel.goSentinelService/GetApplicationSecret" {

		if ctxWithClaims, err = s.authenticate(ctx, ""); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Unauthorized")
		}
	} else if info.FullMethod == "/goSentinel.goSentinelService/CreateApplicationSecret" {
		if ctxWithClaims, err = s.authenticate(ctx, ""); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Unauthorized")
		}
	} else if info.FullMethod == "/goSentinel.goSentinelService/RegisterUser" {
		if ctxWithClaims, err = s.authenticate(ctx, ""); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Unauthorized")
		}
	} else if info.FullMethod == "/goSentinel.goSentinelService/GetUser" {
		request, ok := req.(*goSentinel.GetUserRequest)
		if !ok {
			return nil, status.Errorf(codes.NotFound, "request body not found")
		}
		if ctxWithClaims, err = s.authenticate(ctx, request.UserToken); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid user token")
		}

	} else if info.FullMethod == "/goSentinel.goSentinelService/LoginUser" {
		if ctxWithClaims, err = s.authenticate(ctx, ""); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Unauthorized: %v", err)
		}
	}

	// Calls the handler
	h, err := handler(ctxWithClaims, req)

	// Logging with grpclog (grpclog.LoggerV2)
	grpcLog.Infof("Request - Method:%s\tDuration:%s\tError:%v\n",
		info.FullMethod,
		time.Since(start),
		err)

	return h, err
}

// authorize function authorizes the token received from Metadata

func (s *Service) authenticate(ctx context.Context, userToken string) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token not supplied")
	}

	token := authHeader[0]

	// validateToken function validates the token
	claims, err := s.VerifyJWTToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid application token: %v", err)
	}

	// if claims.ApplicationJWTClaims.ApplicationEmail == nil {
	// 	return nil, status.Errorf(codes.Unauthenticated, "invalid application token: %v", err)
	// }

	if userToken != "" {
		if userClaims, err := s.VerifyJWTToken(userToken); err == nil {
			fmt.Printf("claims are ...%+v", userClaims.UserJWTClaims)
			if userClaims.UserJWTClaims.UserEmail == nil {
				return nil, status.Errorf(codes.Unauthenticated, "invalid user token: %v", err)
			}
			claims.UserJWTClaims = userClaims.UserJWTClaims
		} else {
			return nil, status.Errorf(codes.Unauthenticated, "invalid user token: %v", err)
		}

	}

	ctxWithClaims := context.WithValue(ctx, model.ContextKeyJWTClaims, claims)
	return ctxWithClaims, nil
}
