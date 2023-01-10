package rpc_services

import (
	"GophKeeper/internal/server/app_services/app_service_auth"
	"context"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/model/auth"
	mock "GophKeeper/internal/server/app_services/app_service_auth/mocks"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/auth"
	"GophKeeper/pkg/token"
)

func TestAuthServiceRPC_Login(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authApp := mock.NewMockAuthApp(ctrl)

	tests := []struct {
		name     string
		in       *pb.AuthRequest
		errApp   error
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name: "success",
			in: &pb.AuthRequest{
				Email:    "test@email.com",
				Password: "testPassword",
			},
			wantErr: false,
			errApp:  nil,
		},
		{
			name: "Not found user",
			in: &pb.AuthRequest{
				Email:    "test@email.com",
				Password: "testPassword",
			},
			wantErr:  true,
			wantCode: codes.NotFound,
			errApp:   errs.ErrNotFound,
		},
		{
			name: "Invalid password",
			in: &pb.AuthRequest{
				Email:    "test@email.com",
				Password: "testPassword",
			},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
			errApp:   app_service_auth.ErrInvalidPassword,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var tokenStr string

			if tt.errApp == nil {
				var errJWT error
				tokenStr, errJWT = token.GenerateJWT(tt.in.Email, "")
				require.NoError(t, errJWT)
			}

			cred := auth.Credential{
				Email:    tt.in.Email,
				Password: tt.in.Password,
			}

			authApp.EXPECT().Login(cred).Return(tokenStr, tt.errApp)

			serv := NewAuthServiceRPC(authApp)
			resp, err := serv.Login(context.Background(), tt.in)

			if tt.wantErr {
				if e, ok := status.FromError(err); ok {
					assert.Equal(t, e.Code(), tt.wantCode)
				}
			} else {
				assert.Equal(t, resp.Token, tokenStr)
			}
		})
	}
}

func TestAuthServiceRPC_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authApp := mock.NewMockAuthApp(ctrl)

	tests := []struct {
		name     string
		in       *pb.AuthRequest
		errApp   error
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name: "success",
			in: &pb.AuthRequest{
				Email:    "test@email.com",
				Password: "testPassword",
			},
			wantErr: false,
			errApp:  nil,
		},
		{
			name: "User already exist",
			in: &pb.AuthRequest{
				Email:    "test@email.com",
				Password: "testPassword",
			},
			wantErr:  true,
			wantCode: codes.AlreadyExists,
			errApp:   errs.ErrAlreadyExist,
		},
		{
			name: "Invalid password",
			in: &pb.AuthRequest{
				Email:    "test@email.com",
				Password: "testPassword",
			},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
			errApp:   app_service_auth.ErrInvalidPassword,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var tokenStr string

			if tt.errApp == nil {
				var errJWT error
				tokenStr, errJWT = token.GenerateJWT(tt.in.Email, "")
				require.NoError(t, errJWT)
			}

			cred := auth.Credential{
				Email:    tt.in.Email,
				Password: tt.in.Password,
			}

			authApp.EXPECT().Register(cred).Return(tokenStr, tt.errApp)

			serv := NewAuthServiceRPC(authApp)
			resp, err := serv.Register(context.Background(), tt.in)

			if tt.wantErr {
				if e, ok := status.FromError(err); ok {
					assert.Equal(t, e.Code(), tt.wantCode)
				}
			} else {
				assert.Equal(t, resp.Token, tokenStr)
			}
		})
	}
}

func TestAuthServiceRPC_ChangePassword(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authApp := mock.NewMockAuthApp(ctrl)

	tests := []struct {
		name     string
		email    string
		in       *pb.ChangePasswordRequest
		callApp  bool
		errApp   error
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name:  "Success",
			email: "test@email.com",
			in: &pb.ChangePasswordRequest{
				Password: "qwerty123",
			},
			callApp: true,
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Invalid email",
			in: &pb.ChangePasswordRequest{
				Password: "12345",
			},
			callApp:  false,
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name:  "Invalid password",
			email: "test@email.com",
			in: &pb.ChangePasswordRequest{
				Password: "12345",
			},
			callApp:  true,
			errApp:   app_service_auth.ErrShortPassword,
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name:  "Empty password",
			email: "test@email.com",
			in: &pb.ChangePasswordRequest{
				Password: "",
			},
			callApp:  true,
			errApp:   app_service_auth.ErrShortPassword,
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name:  "Anomaly AppService",
			email: "test@email.com",
			in: &pb.ChangePasswordRequest{
				Password: "qwerty123",
			},
			callApp:  true,
			errApp:   errs.ErrInternal,
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx := context.Background()

			if len(tt.email) != 0 {
				md := metadata.New(map[string]string{"email": tt.email})
				ctx = metadata.NewIncomingContext(ctx, md)
			}

			if tt.callApp {
				if tt.errApp == nil {
					authApp.EXPECT().ChangePassword(tt.email, tt.in.Password).Return(nil)
				} else {
					authApp.EXPECT().ChangePassword(tt.email, tt.in.Password).Return(tt.errApp)
				}
			}

			serv := NewAuthServiceRPC(authApp)
			_, err := serv.ChangePassword(ctx, tt.in)

			if tt.wantErr {
				if e, ok := status.FromError(err); ok {
					assert.Equal(t, e.Code(), tt.wantCode)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
