package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"GophKeeper/pkg/md_ctx"
	"GophKeeper/pkg/token"
)

// ValidateTokenInterceptor - Тест перехватчика для проверки корректности JWT.
func TestValidateTokenInterceptor(t *testing.T) {

	secretKey := ""
	email := "test@email.ru"

	tokenStr, errJWT := token.GenerateJWT(email, "")
	require.NoError(t, errJWT)

	service := func(ctx context.Context, req interface{}) (interface{}, error) {
		emailMD, _ := md_ctx.ValueFromContext(ctx, "email")
		return emailMD, nil
	}

	tests := []struct {
		name      string
		info      *grpc.UnaryServerInfo
		token     string
		wantErr   bool
		wantEmail bool
		wantCode  codes.Code
	}{
		{
			name: "Check unprocessed endpoint Register",
			info: &grpc.UnaryServerInfo{
				FullMethod: "/auth.AuthService/Register",
			},
			wantErr:   false,
			wantEmail: false,
		},
		{
			name: "Check unprocessed endpoint Login",
			info: &grpc.UnaryServerInfo{
				FullMethod: "/auth.AuthService/Login",
			},
			wantErr:   false,
			wantEmail: false,
		},
		{
			name: "Validate valid token",
			info: &grpc.UnaryServerInfo{
				FullMethod: "/auth.AuthService/Any",
			},
			token:   tokenStr,
			wantErr: false,
		},
		{
			name: "Validate empty token",
			info: &grpc.UnaryServerInfo{
				FullMethod: "/auth.AuthService/Any",
			},
			wantErr:  true,
			wantCode: codes.PermissionDenied,
		},
		{
			name: "Validate invalid token",
			info: &grpc.UnaryServerInfo{
				FullMethod: "/auth.AuthService/Any",
			},
			token:    tokenStr + "321",
			wantErr:  true,
			wantCode: codes.PermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			md := metadata.New(map[string]string{"token": tt.token})

			ctx := metadata.NewIncomingContext(context.Background(), md)
			v := ValidateInterceptor{
				secretKey: secretKey,
			}

			emailGet, err := v.ValidateTokenInterceptor(ctx, email, tt.info, service)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				if tt.wantEmail {
					assert.Equal(t, email, emailGet)
				}
			}
		})
	}
}
