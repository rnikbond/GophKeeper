package grpc_service_binary

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/server/model/binary"
	mock "GophKeeper/internal/server/server_grpc/services/grpc_service_binary/mocks"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/binary"
)

func TestBinaryServiceRPC_Create(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	binApp := mock.NewMockBinaryApp(ctrl)

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
				MetaInfo: "desktop.bin",
				Data:     []byte("0101010101"),
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Already exist",
			in: &pb.CreateRequest{
				MetaInfo: "desktop.bin",
				Data:     []byte("0101010101"),
			},
			errApp:   errs.ErrAlreadyExist,
			wantErr:  true,
			wantCode: codes.AlreadyExists,
		},
		{
			name: "Anomaly app service",
			in: &pb.CreateRequest{
				MetaInfo: "desktop.bin",
				Data:     []byte("0101010101"),
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := binary.DataFull{
				MetaInfo: tt.in.MetaInfo,
				Bytes:    tt.in.Data,
			}

			binApp.EXPECT().Create(data).Return(tt.errApp)

			serv := NewBinaryServiceRPC(binApp)
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

func TestBinaryServiceRPC_Change(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	binApp := mock.NewMockBinaryApp(ctrl)

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
				MetaInfo: "desktop.bin",
				Data:     []byte("010101"),
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Not found",
			in: &pb.ChangeRequest{
				MetaInfo: "desktop.bin",
				Data:     []byte("010101"),
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "Anomaly app service",
			in: &pb.ChangeRequest{
				MetaInfo: "desktop.bin",
				Data:     []byte("010101"),
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := binary.DataFull{
				MetaInfo: tt.in.MetaInfo,
				Bytes:    tt.in.Data,
			}

			binApp.EXPECT().Change(data).Return(tt.errApp)

			serv := NewBinaryServiceRPC(binApp)
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

func TestBinaryServiceRPC_Delete(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	binApp := mock.NewMockBinaryApp(ctrl)

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
				MetaInfo: "desktop.bin",
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Not found",
			in: &pb.DeleteRequest{
				MetaInfo: "desktop.bin",
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "Anomaly app service",
			in: &pb.DeleteRequest{
				MetaInfo: "desktop.bin",
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := binary.DataGet{
				MetaInfo: tt.in.MetaInfo,
			}

			binApp.EXPECT().Delete(data).Return(tt.errApp)

			serv := NewBinaryServiceRPC(binApp)
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

func TestBinaryServiceRPC_Get(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	binApp := mock.NewMockBinaryApp(ctrl)

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
				MetaInfo: "desktop.bin",
			},
			out: &pb.GetResponse{
				MetaInfo: "desktop.bin",
				Data:     []byte("010101"),
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Not found",
			in: &pb.GetRequest{
				MetaInfo: "desktop.bin",
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "Anomaly app service",
			in: &pb.GetRequest{
				MetaInfo: "desktop.bin",
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := binary.DataGet{
				MetaInfo: tt.in.MetaInfo,
			}

			outApp := binary.DataFull{
				MetaInfo: tt.in.MetaInfo,
				Bytes:    []byte("010101"),
			}

			binApp.EXPECT().Get(data).Return(outApp, tt.errApp)
			serv := NewBinaryServiceRPC(binApp)
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
