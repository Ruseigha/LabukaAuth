package handler

import (
	"context"
	"errors"

	"github.com/Ruseigha/LabukaAuth/internal/delivery/grpc/proto/proto"
	domainerrors "github.com/Ruseigha/LabukaAuth/internal/domain/errors"
	"github.com/Ruseigha/LabukaAuth/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthHandler implements gRPC AuthServiceServer
type AuthHandler struct {
	proto.UnimplementedAuthServiceServer // Forward compatibility
	authService                          usecase.AuthUseCase
}

// NewAuthHandler creates a new gRPC auth handler
func NewAuthHandler(authService usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Signup implements gRPC Signup RPC
func (h *AuthHandler) Signup(ctx context.Context, req *proto.SignupRequest) (*proto.AuthResponse, error) {
	// Validate request
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Call use case
	resp, err := h.authService.Signup(ctx, usecase.SignupRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	// Return response
	return &proto.AuthResponse{
		UserId:       resp.UserID,
		Email:        resp.Email,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// Login implements gRPC Login RPC
func (h *AuthHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {
	// Validate
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Call use case
	resp, err := h.authService.Login(ctx, usecase.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &proto.AuthResponse{
		UserId:       resp.UserID,
		Email:        resp.Email,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// RefreshToken implements gRPC RefreshToken RPC
func (h *AuthHandler) RefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (*proto.AuthResponse, error) {
	// Validate
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	// Call use case
	resp, err := h.authService.RefreshToken(ctx, req.RefreshToken)

	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &proto.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// ValidateToken implements gRPC ValidateToken RPC
func (h *AuthHandler) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	// Validate
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access_token is required")
	}

	// Call use case
	claims, err := h.authService.ValidateToken(ctx, req.AccessToken)

	if err != nil {
		return nil, mapDomainErrorToGRPC(err)
	}

	return &proto.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Email:  claims.Email,
	}, nil
}

// mapDomainErrorToGRPC maps domain errors to gRPC status codes
func mapDomainErrorToGRPC(err error) error {
	// Map domain errors to gRPC codes
	switch {
	case errors.Is(err, domainerrors.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, domainerrors.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, err.Error())

	case errors.Is(err, domainerrors.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())

	case errors.Is(err, domainerrors.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, domainerrors.ErrConflict):
		return status.Error(codes.AlreadyExists, err.Error())

	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
