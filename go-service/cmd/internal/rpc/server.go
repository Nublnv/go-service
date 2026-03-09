package rpc

import (
	"net"

	"github.com/Nublnv/go-service/cmd/internal/db"
	"github.com/Nublnv/go-service/cmd/internal/errorHandler"
	"github.com/Nublnv/go-service/cmd/internal/middleware"
	v1 "github.com/Nublnv/go-service/cmd/internal/rpc/payslip/v1"
	payslipv1 "github.com/Nublnv/go-service/proto/payslip/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GetRpcServer(pg *pgxpool.Pool, click *db.Pool, tlsCert string, tlsKey string) (*grpc.Server, error) {
	creds, err := credentials.NewServerTLSFromFile(tlsCert, tlsKey)
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			errorHandler.RpcWrap,
			middleware.AuthMiddlewareRpc,
		),
		grpc.Creds(creds),
	)

	svc := &v1.DocumentService{
		Pg:    pg,
		Click: click,
	}

	payslipv1.RegisterDocumentServiceServer(server, svc)

	return server, nil
}

func ServeServer(svc *grpc.Server, l net.Listener) func() {
	return func() {
		if err := svc.Serve(l); err != nil {
		}
	}
}
