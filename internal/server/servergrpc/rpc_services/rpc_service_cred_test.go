package rpc_services

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/model/cred"
	mock "GophKeeper/internal/server/app_services/app_service_credential/mocks"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/data/credential"
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
				Email:    "test@email.com",
				MetaInfo: "www.test.ru",
				Password: "testPwd",
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Already exist",
			in: &pb.CreateRequest{
				Email:    "test@email.com",
				MetaInfo: "www.test.ru",
				Password: "testPwd",
			},
			errApp:   errs.ErrAlreadyExist,
			wantErr:  true,
			wantCode: codes.AlreadyExists,
		},
		{
			name: "Anomaly app service",
			in: &pb.CreateRequest{
				Email:    "test@email.com",
				MetaInfo: "www.test.ru",
				Password: "testPwd",
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := cred.CredentialFull{
				Email:    tt.in.Email,
				MetaInfo: tt.in.MetaInfo,
				Password: tt.in.Password,
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
				Email:    "test@email.com",
				MetaInfo: "www.test.ru",
				Password: "testPwd",
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Not found",
			in: &pb.ChangeRequest{
				Email:    "test@email.com",
				MetaInfo: "www.test.ru",
				Password: "testPwd",
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "Anomaly app service",
			in: &pb.ChangeRequest{
				Email:    "test@email.com",
				MetaInfo: "www.test.ru",
				Password: "testPwd",
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := cred.CredentialFull{
				Email:    tt.in.Email,
				MetaInfo: tt.in.MetaInfo,
				Password: tt.in.Password,
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
				Email:    "test@email.com",
				MetaInfo: "www.test.ru",
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Not found",
			in: &pb.DeleteRequest{
				Email:    "test@email.com",
				MetaInfo: "www.test.ru",
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "Anomaly app service",
			in: &pb.DeleteRequest{
				Email:    "test@email.com",
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
				Email:    tt.in.Email,
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
				Email:    "test@email.com",
				MetaInfo: "www.test.ru",
			},
			out: &pb.GetResponse{
				Email:    "test@email.com",
				MetaInfo: "www.test.ru",
				Password: "testPwd",
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Not found",
			in: &pb.GetRequest{
				Email:    "test@email.com",
				MetaInfo: "www.test.ru",
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "Anomaly app service",
			in: &pb.GetRequest{
				Email:    "test@email.com",
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
				Email:    tt.in.Email,
				MetaInfo: tt.in.MetaInfo,
			}

			outApp := cred.CredentialFull{
				Email:    tt.in.Email,
				MetaInfo: tt.in.MetaInfo,
				Password: "testPwd",
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
