package grpc_service_card

import (
	"GophKeeper/internal/client/grpc_services/grpc_service_card"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/model/card"
	mock "GophKeeper/internal/server/server_grpc/services/grpc_service_card/mocks"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/card"
)

func TestCardServiceRPC_Create(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cardApp := mock.NewMockCardApp(ctrl)

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
				MetaInfo: "MirPay",
				Number:   []byte("464289760410976"),
				Period:   []byte("10.2030"),
				CVV:      []byte("111"),
				FullName: []byte("Test Test"),
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Check already exist",
			in: &pb.CreateRequest{
				MetaInfo: "MirPay",
				Number:   []byte("4648289760410976"),
				Period:   []byte("10.2030"),
				CVV:      []byte("111"),
				FullName: []byte("Test Test"),
			},
			errApp:   errs.ErrAlreadyExist,
			wantErr:  true,
			wantCode: codes.AlreadyExists,
		},
		{
			name: "Check invalid number",
			in: &pb.CreateRequest{
				MetaInfo: "MirPay",
				Number:   []byte("464828760410976"),
				Period:   []byte("10.2030"),
				CVV:      []byte("111"),
				FullName: []byte("Test Test"),
			},
			errApp:   grpc_service_card.ErrInvalidNumber,
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Check invalid period",
			in: &pb.CreateRequest{
				MetaInfo: "MirPay",
				Number:   []byte("4648289760410976"),
				Period:   []byte("10/2030"),
				CVV:      []byte("111"),
				FullName: []byte("Test Test"),
			},
			errApp:   grpc_service_card.ErrInvalidPeriod,
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Check invalid CVV: chars",
			in: &pb.CreateRequest{
				MetaInfo: "MirPay",
				Number:   []byte("4648289760410976"),
				Period:   []byte("10.2030"),
				CVV:      []byte("aaa"),
				FullName: []byte("Test Test"),
			},
			errApp:   grpc_service_card.ErrInvalidCVV,
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Check invalid CVV: short",
			in: &pb.CreateRequest{
				MetaInfo: "MirPay",
				Number:   []byte("4648289760410976"),
				Period:   []byte("10.2030"),
				CVV:      []byte("11"),
				FullName: []byte("Test Test"),
			},
			errApp:   grpc_service_card.ErrInvalidCVV,
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Check invalid full name",
			in: &pb.CreateRequest{
				MetaInfo: "MirPay",
				Number:   []byte("4648289760410976"),
				Period:   []byte("10.2030"),
				CVV:      []byte("111"),
				FullName: []byte("Te"),
			},
			errApp:   grpc_service_card.ErrInvalidFullName,
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := card.DataCardFull{
				MetaInfo: tt.in.MetaInfo,
				Number:   string(tt.in.Number),
				Period:   string(tt.in.Period),
				CVV:      string(tt.in.CVV),
				FullName: string(tt.in.FullName),
			}

			cardApp.EXPECT().Create(data).Return(tt.errApp)

			serv := NewCardServiceRPC(cardApp)
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

func TestCardServiceRPC_Change(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cardApp := mock.NewMockCardApp(ctrl)

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
				MetaInfo: "MirPay",
				Number:   []byte("464289760410976"),
				Period:   []byte("10.2030"),
				CVV:      []byte("111"),
				FullName: []byte("Test Test"),
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Check not found",
			in: &pb.ChangeRequest{
				MetaInfo: "MirPay",
				Number:   []byte("4648289760410976"),
				Period:   []byte("10.2030"),
				CVV:      []byte("111"),
				FullName: []byte("Test Test"),
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
		{
			name: "Check invalid number",
			in: &pb.ChangeRequest{
				MetaInfo: "MirPay",
				Number:   []byte("464828760410976"),
				Period:   []byte("10.2030"),
				CVV:      []byte("111"),
				FullName: []byte("Test Test"),
			},
			errApp:   grpc_service_card.ErrInvalidNumber,
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Check invalid period",
			in: &pb.ChangeRequest{
				MetaInfo: "MirPay",
				Number:   []byte("4648289760410976"),
				Period:   []byte("10/2030"),
				CVV:      []byte("111"),
				FullName: []byte("Test Test"),
			},
			errApp:   grpc_service_card.ErrInvalidPeriod,
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Check invalid CVV: chars",
			in: &pb.ChangeRequest{
				MetaInfo: "MirPay",
				Number:   []byte("4648289760410976"),
				Period:   []byte("10.2030"),
				CVV:      []byte("aaa"),
				FullName: []byte("Test Test"),
			},
			errApp:   grpc_service_card.ErrInvalidCVV,
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Check invalid CVV: short",
			in: &pb.ChangeRequest{
				MetaInfo: "MirPay",
				Number:   []byte("4648289760410976"),
				Period:   []byte("10.2030"),
				CVV:      []byte("11"),
				FullName: []byte("Test Test"),
			},
			errApp:   grpc_service_card.ErrInvalidCVV,
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Check invalid full name",
			in: &pb.ChangeRequest{
				MetaInfo: "MirPay",
				Number:   []byte("4648289760410976"),
				Period:   []byte("10.2030"),
				CVV:      []byte("111"),
				FullName: []byte("Te"),
			},
			errApp:   grpc_service_card.ErrInvalidFullName,
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := card.DataCardFull{
				MetaInfo: tt.in.MetaInfo,
				Number:   string(tt.in.Number),
				Period:   string(tt.in.Period),
				CVV:      string(tt.in.CVV),
				FullName: string(tt.in.FullName),
			}

			cardApp.EXPECT().Change(data).Return(tt.errApp)

			serv := NewCardServiceRPC(cardApp)
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

func TestCardServiceRPC_Delete(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cardApp := mock.NewMockCardApp(ctrl)

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
				MetaInfo: "MirPay",
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Check not found",
			in: &pb.DeleteRequest{
				MetaInfo: "MirPay",
			},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := card.DataCardGet{
				MetaInfo: tt.in.MetaInfo,
			}

			cardApp.EXPECT().Delete(data).Return(tt.errApp)

			serv := NewCardServiceRPC(cardApp)
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

func TestCardServiceRPC_Get(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cardApp := mock.NewMockCardApp(ctrl)

	tests := []struct {
		name     string
		in       *pb.GetRequest
		out      *pb.GetResponse
		outApp   card.DataCardFull
		errApp   error
		wantErr  bool
		wantCode codes.Code
	}{
		{
			name: "Success",
			in: &pb.GetRequest{
				MetaInfo: "MirPay",
			},
			out: &pb.GetResponse{
				Number:   []byte("4648289760410976"),
				Period:   []byte("10.2030"),
				CVV:      []byte("111"),
				FullName: []byte("Test Test"),
			},
			outApp: card.DataCardFull{
				Number:   "4648289760410976",
				Period:   "10.2030",
				CVV:      "111",
				FullName: "Test Test",
			},
			errApp:  nil,
			wantErr: false,
		},
		{
			name: "Check not found",
			in: &pb.GetRequest{
				MetaInfo: "MirPay",
			},
			out:      &pb.GetResponse{},
			outApp:   card.DataCardFull{},
			errApp:   errs.ErrNotFound,
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := card.DataCardGet{
				MetaInfo: tt.in.MetaInfo,
			}

			cardApp.EXPECT().Get(data).Return(tt.outApp, tt.errApp)

			serv := NewCardServiceRPC(cardApp)
			out, err := serv.Get(context.Background(), tt.in)

			if tt.wantErr {
				if e, ok := status.FromError(err); ok {
					assert.Equal(t, e.Code(), tt.wantCode)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, out, tt.out)
			}
		})
	}
}
