package auth

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
	svcerrors "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/errors"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

type AuthServerImpl struct {
	authpb.UnimplementedAuthServiceServer
	authUC AuthUseCase
}

func NewAuthServer(authUC AuthUseCase) *AuthServerImpl {
	return &AuthServerImpl{
		authUC: authUC,
	}
}

func (s *AuthServerImpl) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	regReq := RegisterToRequest(req)
	user, err := s.authUC.Register(ctx, regReq)
	if err != nil {
		if errors.Is(err, svcerrors.ErrLoginExists) {
			return nil, status.Error(codes.AlreadyExists, string(models.ErrCodeLoginExists))
		}
		if errors.Is(err, svcerrors.ErrEmailExists) {
			return nil, status.Error(codes.AlreadyExists, string(models.ErrCodeEmailExists))
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return user, nil
}

func (s *AuthServerImpl) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	logReq := LoginToRequest(req)
	user, err := s.authUC.Login(ctx, logReq)
	if err != nil {
		if errors.Is(err, svcerrors.ErrInvalidCredentials) ||  errors.Is(err, svcerrors.ErrUserNotFound) {
			return nil, status.Error(codes.Unauthenticated, string(models.ErrCodeInvalidCredentials))
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}

	return user, nil
}

func (s *AuthServerImpl) UpdateProfile(ctx context.Context, req *authpb.UpdateProfileRequest) (*authpb.ProfileResponse, error) {
	profileReq := UpdateProfileToRequest(req)
	profile, err := s.authUC.UpdateProfile(ctx, profileReq)
	if err != nil {
		if errors.Is(err, svcerrors.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, string(models.ErrCodeUserNotFound))
		}
		if errors.Is(err, svcerrors.ErrLoginExists) {
			return nil, status.Error(codes.AlreadyExists, string(models.ErrCodeEmailExists))
		}
		if errors.Is(err, svcerrors.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, string(models.ErrCodeForbidden	))
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return profile, err
}

func (s *AuthServerImpl) GetProfile(ctx context.Context, userID *authpb.UserID) (*authpb.ProfileResponse, error) {
	id := ProtoIDtoInt(userID)
	profile, err := s.authUC.GetProfile(ctx, id)
	if err != nil {
		if errors.Is(err, svcerrors.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, string(models.ErrCodeUserNotFound))
		}
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return profile, err
}

func (s *AuthServerImpl) GetCSRF(ctx context.Context, _ *emptypb.Empty) (*authpb.CSRFTokenResponse, error) {
	token, err := s.authUC.GetCSRFToken(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, string(models.ErrCodeInternalError))
	}
	return token, err
}


