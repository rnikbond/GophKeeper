//go:generate mockgen -source grpc_service_cred.go -destination mocks/grpc_service_cred_mock.go -package grpc_service_cred
package grpc_service_cred

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/client/model/cred_model"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/credential"
)

type CredService struct {
	rpc    pb.CredentialServiceClient
	logger *zap.Logger
}

// NewService - Создание экземпляра сервиса для текстовых данных.
func NewService(conn *grpc.ClientConn) *CredService {
	return &CredService{
		rpc:    pb.NewCredentialServiceClient(conn),
		logger: zap.L(),
	}
}

func (serv CredService) Create(data cred_model.Credential, token string) error {

	dataReq := &pb.CreateRequest{
		MetaInfo: data.MetaInfo,
		Email:    data.Login,
		Password: data.Password,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if _, err := serv.rpc.Create(ctx, dataReq); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				return errs.ErrAlreadyExist

			default:
				serv.logger.Error("unknown gRPC error in cred service Create()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}

		return errs.ErrInternal
	}

	return nil
}

func (serv CredService) Get(meta string, token string) (cred_model.Credential, error) {
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
				return cred_model.Credential{}, errs.ErrNotFound

			default:
				serv.logger.Error("unknown gRPC error in cred service Get()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}

		return cred_model.Credential{}, errs.ErrInternal
	}

	return cred_model.Credential{
		MetaInfo: meta,
		Login:    resp.Email,
		Password: resp.Password,
	}, nil
}

func (serv CredService) Delete(meta string, token string) error {
	data := &pb.DeleteRequest{
		MetaInfo: meta,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := serv.rpc.Delete(ctx, data)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return errs.ErrNotFound

			default:
				serv.logger.Error("unknown gRPC error in cred service Delete()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}
		return errs.ErrInternal
	}

	return nil
}

func (serv CredService) Change(data cred_model.Credential, token string) error {

	dataReq := &pb.ChangeRequest{
		MetaInfo: data.MetaInfo,
		Email:    data.Login,
		Password: data.Password,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if _, err := serv.rpc.Change(ctx, dataReq); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return errs.ErrNotFound

			default:
				serv.logger.Error("unknown gRPC error in cred service Change()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}
		return errs.ErrInternal
	}

	return nil
}
