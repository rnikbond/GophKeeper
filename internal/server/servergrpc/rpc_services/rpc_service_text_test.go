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

	"GophKeeper/internal/model/text"
	mock "GophKeeper/internal/server/app_services/app_service_text/mocks"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/data/text"
)

func TestTextServiceRPC_Create(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	textApp := mock.NewMockTextApp(ctrl)

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
				MetaInfo: "book1",
				Text:     "testText",
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Already exist",
			in: &pb.CreateRequest{
				MetaInfo: "book1",
				Text:     "testText",
			},
			errApp:   errs.ErrAlreadyExist,
			wantErr:  true,
			wantCode: codes.AlreadyExists,
		},
		{
			name: "Anomaly app service",
			in: &pb.CreateRequest{
				MetaInfo: "book1",
				Text:     "testText",
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := text.DataTextFull{
				MetaInfo: tt.in.MetaInfo,
				Text:     tt.in.Text,
			}

			textApp.EXPECT().Create(data).Return(tt.errApp)

			serv := NewTextServiceRPC(textApp)
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

func TestTextServiceRPC_Change(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	textApp := mock.NewMockTextApp(ctrl)

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
				MetaInfo: "book1",
				Text:     "testText",
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Not found",
			in: &pb.ChangeRequest{
				MetaInfo: "book1",
				Text:     "testText",
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "Anomaly app service",
			in: &pb.ChangeRequest{
				MetaInfo: "book1",
				Text:     "testText",
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := text.DataTextFull{
				MetaInfo: tt.in.MetaInfo,
				Text:     tt.in.Text,
			}

			textApp.EXPECT().Change(data).Return(tt.errApp)

			serv := NewTextServiceRPC(textApp)
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

func TestTextServiceRPC_Delete(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	textApp := mock.NewMockTextApp(ctrl)

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
				MetaInfo: "book1",
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Not found",
			in: &pb.DeleteRequest{
				MetaInfo: "book1",
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "Anomaly app service",
			in: &pb.DeleteRequest{
				MetaInfo: "book1",
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := text.DataTextGet{
				MetaInfo: tt.in.MetaInfo,
			}

			textApp.EXPECT().Delete(data).Return(tt.errApp)

			serv := NewTextServiceRPC(textApp)
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

func TestTextServiceRPC_Get(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	textApp := mock.NewMockTextApp(ctrl)

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
				MetaInfo: "book1",
			},
			out: &pb.GetResponse{
				MetaInfo: "book1",
				Text:     "testText",
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Not found",
			in: &pb.GetRequest{
				MetaInfo: "book1",
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "Anomaly app service",
			in: &pb.GetRequest{
				MetaInfo: "book1",
			},
			errApp:   fmt.Errorf("unknown error"),
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := text.DataTextGet{
				MetaInfo: tt.in.MetaInfo,
			}

			outApp := text.DataTextFull{
				MetaInfo: tt.in.MetaInfo,
				Text:     "testText",
			}

			textApp.EXPECT().Get(data).Return(outApp, tt.errApp)
			serv := NewTextServiceRPC(textApp)
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
