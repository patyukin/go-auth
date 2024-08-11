package auth

import (
	"context"

	authpb "github.com/bogatyr285/auth-go/pkg/server/grpc/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandlers struct {
	authpb.UnimplementedAuthServiceServer
}

func NewAuthHandlers() *AuthHandlers {
	return &AuthHandlers{}
}

func (h *AuthHandlers) RegisterUser(context.Context, *authpb.RegisterUserRequest) (*authpb.RegisterUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (h *AuthHandlers) LoginUser(context.Context, *authpb.LoginUserRequest) (*authpb.LoginUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (h *AuthHandlers) UserInfo(context.Context, *authpb.UserInfoRequest) (*authpb.UserInfoResponse, error) {
	return &authpb.UserInfoResponse{
		User: &authpb.User{
			UserId: uuid.New().String(),
			ContactMethod: &authpb.User_PhoneNumber{
				PhoneNumber: "7-800-555-3535",
			},
		},
	}, nil
	//return nil, status.Error(codes.Unimplemented, "")
}
