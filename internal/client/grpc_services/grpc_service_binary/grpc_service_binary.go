package grpc_service_binary

import (
	"GophKeeper/internal/client/model/binary_model"
	"context"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/binary"
)

type BinaryOptions func(c *BinaryService)

type BinaryService struct {
	rpc    pb.BinaryServiceClient
	logger *zap.Logger
}

// NewService - Создание экземпляра сервиса для текстовых данных.
func NewService(conn *grpc.ClientConn) *BinaryService {
	return &BinaryService{
		rpc:    pb.NewBinaryServiceClient(conn),
		logger: zap.L(),
	}

}

func (serv BinaryService) Create(data binary_model.Binary, token string) error {
	dataReq := &pb.CreateRequest{
		MetaInfo: data.MetaInfo,
		Data:     data.Data,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if _, err := serv.rpc.Create(ctx, dataReq); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				return errs.ErrAlreadyExist

			default:
				if strings.Contains(err.Error(), "larger than max") {
					return errs.ErrLargeData
				}

				serv.logger.Error("unknown gRPC error in binary service Create()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}
		return errs.ErrInternal
	}

	return nil
}

func (serv BinaryService) Get(meta, token string) (binary_model.Binary, error) {
	data := &pb.GetRequest{
		MetaInfo: meta,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := serv.rpc.Get(ctx, data)
	if err != nil {
		if e, ok := status.FromError(err); ok {

			switch e.Code() {
			case codes.NotFound:
				return binary_model.Binary{}, errs.ErrNotFound

			default:
				serv.logger.Error("unknown gRPC error in binary service Get()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}
		return binary_model.Binary{}, errs.ErrInternal
	}

	return binary_model.Binary{
		MetaInfo: meta,
		Data:     resp.Data,
	}, nil
}

func (serv BinaryService) Delete(meta, token string) error {
	data := &pb.DeleteRequest{
		MetaInfo: meta,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if _, err := serv.rpc.Delete(ctx, data); err != nil {
		if e, ok := status.FromError(err); ok {

			switch e.Code() {
			case codes.NotFound:
				return errs.ErrNotFound

			default:
				serv.logger.Error("unknown gRPC error in binary service Delete()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}
		return errs.ErrInternal
	}

	return nil
}

func (serv BinaryService) Change(data binary_model.Binary, token string) error {
	dataReq := &pb.ChangeRequest{
		MetaInfo: data.MetaInfo,
		Data:     data.Data,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := serv.rpc.Change(ctx, dataReq)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return errs.ErrNotFound

			default:
				if strings.Contains(err.Error(), "larger than max") {
					return errs.ErrLargeData
				}

				serv.logger.Error("unknown gRPC error in binary service Change()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}
		return errs.ErrInternal
	}

	return nil
}
