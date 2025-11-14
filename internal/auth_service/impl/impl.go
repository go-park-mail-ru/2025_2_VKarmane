package auth

import (
	"context"

	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
	"google.golang.org/protobuf/types/known/emptypb"
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
		return nil, err
	}
	return user, nil
}

func (s *AuthServerImpl) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	logReq := LoginToRequest(req)
	user, err := s.authUC.Login(ctx, logReq)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthServerImpl) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	empty, err := s.authUC.Logout(ctx)
	return empty, err
}

func (s *AuthServerImpl) UpdateProfile(ctx context.Context, req *authpb.UpdateProfileRequest) (*authpb.ProfileResponse, error) {
	profileReq := UpdateProfileToRequest(req)
	profile, err := s.authUC.UpdateProfile(ctx, profileReq)
	if err != nil {
		return nil, err
	}
	return profile, err
}

func (s *AuthServerImpl) GetProfile(ctx context.Context, userID *authpb.UserID) (*authpb.ProfileResponse, error) {
	id := ProtoIDtoInt(userID)
	profile, err := s.authUC.GetProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	return profile, err
}

func (s *AuthServerImpl) GetCSRF(ctx context.Context, _ *emptypb.Empty) (*authpb.CSRFTokenResponse, error) {
	token, err := s.authUC.GetCSRFToken(ctx)
	if err != nil {
		return nil, err
	}
	return token, err
}


