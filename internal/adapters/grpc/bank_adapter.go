package grpc

import (
	"context"
	"time"

	"github.com/VallabhSLEPAM/go-with-grpc/protogen/go/bank"
	"google.golang.org/genproto/googleapis/type/date"
)

func (adapter *GRPCAdapter) GetCurrentBalance(ctx context.Context, req *bank.CurrentBalanceRequest) (*bank.CurrentBalanceResponse, error) {

	now := time.Now()
	bal := adapter.bankService.FindCurrentBalance(req.AccountNumber)

	return &bank.CurrentBalanceResponse{
		Amount: bal,
		CurrentDate: &date.Date{
			Year:  int32(now.Year()),
			Month: int32(now.Month()),
			Day:   int32(now.Day()),
		}}, nil
}
