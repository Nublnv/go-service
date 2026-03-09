package v1

import (
	"context"

	"github.com/Nublnv/go-service/cmd/internal/db"
	payslipv1 "github.com/Nublnv/go-service/proto/payslip/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DocumentService struct {
	Pg    *pgxpool.Pool
	Click *db.Pool
	payslipv1.UnimplementedDocumentServiceServer
}

func (*DocumentService) GeneratePayslip(
	ctx context.Context,
	in *payslipv1.GeneratePayslipRequest,
) (*payslipv1.GeneratePayslipResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Methon not implemented")
}

func (*DocumentService) GetPayslip(context.Context, *payslipv1.GetPayslipRequest) (*payslipv1.GetPayslipResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Methon not implemented")
}

// Method uses for send email with payslip
func (s *DocumentService) SendPayslip(ctx context.Context, r *payslipv1.SendPaylipRequest) (*payslipv1.SendPayslipResponse, error) {
	pg, err := s.Pg.Acquire(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Cannot connect to postgres")
	}

	var email string

	err = pg.QueryRow(ctx, "SELECT email FROM auth.users WHERE user_id = $1", r.Userid).Scan(&email)
	if err != nil || email == "" {
		return nil, status.Error(codes.Internal, "Cannot find user email")
	}

	// TODO: add email sending implementation

	return &payslipv1.SendPayslipResponse{
		Ok: true,
	}, nil
}
