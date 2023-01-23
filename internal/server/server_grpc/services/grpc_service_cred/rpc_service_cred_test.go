package grpc_service_cred

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/server/model/cred"
	mock "GophKeeper/internal/server/server_grpc/services/grpc_service_cred/mocks"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/credential"
)

func TestCredServiceRPC_Create(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	credApp := mock.NewMockCredentialApp(ctrl)

	tests := []struct {
		name     string
		in       *pb.CreateRequest
		errApp   error
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name: "Success",
			in: &pb.CreateRequest{
				MetaInfo: "www.test.ru",
				Email:    []byte("test@email.com"),
				Password: []byte("testPwd"),
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Already exist",
			in: &pb.CreateRequest{
				MetaInfo: "www.test.ru",
				Email:    []byte("test@email.com"),
				Password: []byte("testPwd"),
			},
			errApp:   errs.ErrAlreadyExist,
			wantErr:  true,
			wantCode: codes.AlreadyExists,
		},
		{
			name: "Anomaly app service",
			in: &pb.CreateRequest{
				MetaInfo: "www.test.ru",
				Email:    []byte("test@email.com"),
				Password: []byte("testPwd"),
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := cred.CredentialFull{
				MetaInfo: tt.in.MetaInfo,
				Email:    string(tt.in.Email),
				Password: string(tt.in.Password),
			}

			credApp.EXPECT().Create(data).Return(tt.errApp)

			serv := NewCredServiceRPC(credApp)
			_, err := serv.Create(context.Background(), tt.in)

			if tt.wantErr {
				if e, ok := status.FromError(err); ok {
					assert.Equal(t, e.Code(), tt.wantCode)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCredServiceRPC_Change(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	credApp := mock.NewMockCredentialApp(ctrl)

	tests := []struct {
		name     string
		in       *pb.ChangeRequest
		errApp   error
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name: "Success",
			in: &pb.ChangeRequest{
				MetaInfo: "www.test.ru",
				Email:    []byte("test@email.com"),
				Password: []byte("testPwd"),
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Not found",
			in: &pb.ChangeRequest{
				MetaInfo: "www.test.ru",
				Email:    []byte("test@email.com"),
				Password: []byte("testPwd"),
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "Anomaly app service",
			in: &pb.ChangeRequest{
				MetaInfo: "www.test.ru",
				Email:    []byte("test@email.com"),
				Password: []byte("testPwd"),
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := cred.CredentialFull{
				MetaInfo: tt.in.MetaInfo,
				Email:    string(tt.in.Email),
				Password: string(tt.in.Password),
			}

			credApp.EXPECT().Change(data).Return(tt.errApp)

			serv := NewCredServiceRPC(credApp)
			_, err := serv.Change(context.Background(), tt.in)

			if tt.wantErr {
				if e, ok := status.FromError(err); ok {
					assert.Equal(t, e.Code(), tt.wantCode)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCredServiceRPC_Delete(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	credApp := mock.NewMockCredentialApp(ctrl)

	tests := []struct {
		name     string
		in       *pb.DeleteRequest
		errApp   error
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name: "Success",
			in: &pb.DeleteRequest{
				MetaInfo: "www.test.ru",
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Not found",
			in: &pb.DeleteRequest{
				MetaInfo: "www.test.ru",
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "Anomaly app service",
			in: &pb.DeleteRequest{
				MetaInfo: "www.test.ru",
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := cred.CredentialGet{
				MetaInfo: tt.in.MetaInfo,
			}

			credApp.EXPECT().Delete(data).Return(tt.errApp)

			serv := NewCredServiceRPC(credApp)
			_, err := serv.Delete(context.Background(), tt.in)

			if tt.wantErr {
				if e, ok := status.FromError(err); ok {
					assert.Equal(t, e.Code(), tt.wantCode)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCredServiceRPC_Get(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	credApp := mock.NewMockCredentialApp(ctrl)

	tests := []struct {
		name     string
		in       *pb.GetRequest
		out      *pb.GetResponse
		errApp   error
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name: "Success",
			in: &pb.GetRequest{
				MetaInfo: "www.test.ru",
			},
			out: &pb.GetResponse{
				Email:    []byte("test@email.com"),
				Password: []byte("testPwd"),
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Not found",
			in: &pb.GetRequest{
				MetaInfo: "www.test.ru",
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "Anomaly app service",
			in: &pb.GetRequest{
				MetaInfo: "www.test.ru",
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := cred.CredentialGet{
				MetaInfo: tt.in.MetaInfo,
			}

			var outApp cred.CredentialFull

			if tt.out != nil {
				outApp = cred.CredentialFull{
					MetaInfo: tt.in.MetaInfo,
					Email:    string(tt.out.Email),
					Password: string(tt.out.Password),
				}
			}

			credApp.EXPECT().Get(data).Return(outApp, tt.errApp)

			serv := NewCredServiceRPC(credApp)
			get, err := serv.Get(context.Background(), tt.in)

			if tt.wantErr {
				if e, ok := status.FromError(err); ok {
					assert.Equal(t, e.Code(), tt.wantCode)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, get, tt.out)
			}
		})
	}
}
