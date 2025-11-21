package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	svcerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/errors"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type AuthServiceServer struct {
	authpb.UnimplementedAuthServiceServer
	authUC AuthUseCase
}

func NewAuthServer(authUC AuthUseCase) *AuthServiceServer {
	return &AuthServiceServer{
		authUC: authUC,
	}
}

func (s *AuthServiceServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	regReq := RegisterToRequest(req)
	user, err := s.authUC.Register(ctx, regReq)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range svcerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to create user", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to create user, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))

	}
	return user, nil
}

func (s *AuthServiceServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	logReq := LoginToRequest(req)
	user, err := s.authUC.Login(ctx, logReq)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range svcerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to login user", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to login user, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}

	return user, nil
}

func (s *AuthServiceServer) UpdateProfile(ctx context.Context, req *authpb.UpdateProfileRequest) (*authpb.ProfileResponse, error) {
	profileReq := UpdateProfileToRequest(req)
	profile, err := s.authUC.UpdateProfile(ctx, profileReq)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range svcerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to update user", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to update user, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return profile, err
}

func (s *AuthServiceServer) GetProfile(ctx context.Context, userID *authpb.UserID) (*authpb.ProfileResponse, error) {
	id := ProtoIDtoInt(userID)
	profile, err := s.authUC.GetProfile(ctx, id)
	if err != nil {
		logger := logger.FromContext(ctx)
		for targetErr, resp := range svcerrors.ErrorMap {
			if errors.Is(err, targetErr) {
				if logger != nil {
					logger.Error("Failed to get user", "error", err)
				}
				return nil, status.Error(resp.Code, resp.Msg)
			}
		}
		if logger != nil {
			logger.Error("Failed to get user, internal error", "error", err)
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return profile, err
}

func (s *AuthServiceServer) GetCSRF(ctx context.Context, _ *emptypb.Empty) (*authpb.CSRFTokenResponse, error) {
	token, err := s.authUC.GetCSRFToken(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return token, err
}
