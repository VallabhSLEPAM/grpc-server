package grpc

import (
	"context"
	"io"
	"log"
	"time"

	protobank "github.com/VallabhSLEPAM/go-with-grpc/protogen/go/bank"
	"github.com/VallabhSLEPAM/grpc-server/internal/application/domain/bank"

	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/genproto/googleapis/type/datetime"

	"google.golang.org/grpc"
)

func (adapter *GRPCAdapter) GetCurrentBalance(ctx context.Context, req *protobank.CurrentBalanceRequest) (*protobank.CurrentBalanceResponse, error) {

	now := time.Now()
	bal := adapter.bankService.FindCurrentBalance(req.AccountNumber)

	return &protobank.CurrentBalanceResponse{
		Amount: bal,
		CurrentDate: &date.Date{
			Year:  int32(now.Year()),
			Month: int32(now.Month()),
			Day:   int32(now.Day()),
		}}, nil
}

func (adapter *GRPCAdapter) FetchExchangeRates(exchangeRateReq *protobank.ExchangeRateRequest, stream grpc.ServerStreamingServer[protobank.ExchangeRateResponse]) error {

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			log.Println("Client cancelled the request. Exiting ...")
			return nil
		default:
			now := time.Now().Truncate(time.Second)
			rate := adapter.bankService.FindExchangeRate(exchangeRateReq.FromCurrency, exchangeRateReq.ToCurrency, now)

			stream.Send(
				&protobank.ExchangeRateResponse{
					FromCurrency: exchangeRateReq.FromCurrency,
					ToCurrency:   exchangeRateReq.ToCurrency,
					Rate:         rate,
					Timestamp:    now.Format(time.RFC3339),
				},
			)
			log.Printf("Exchange rate sent to client from %v to %v: %v", exchangeRateReq.FromCurrency, exchangeRateReq.ToCurrency, rate)
			time.Sleep(3 * time.Second)
		}
	}

}

func toTime(dt *datetime.DateTime) (time.Time, error) {

	if dt == nil {
		now := time.Now()
		dt = &datetime.DateTime{
			Year:    int32(now.Year()),
			Month:   int32(now.Month()),
			Day:     int32(now.Day()),
			Hours:   int32(now.Hour()),
			Minutes: int32(now.Minute()),
			Seconds: int32(now.Second()),
			Nanos:   int32(now.Nanosecond()),
		}
	}
	res := time.Date(int(dt.Year), time.Month(dt.Month), int(dt.Day), int(dt.Hours), int(dt.Minutes), int(dt.Seconds), int(dt.Nanos), time.UTC)
	return res, nil
}

func (grpcAdapter *GRPCAdapter) SummarizeTransactions(stream grpc.ClientStreamingServer[protobank.Transaction, protobank.TransactionSummary]) error {

	tsummary := bank.TransactionSummary{
		SumIn:    0,
		SumOut:   0,
		SumTotal: 0,
	}

	acct := ""
	now := time.Now()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := protobank.TransactionSummary{
				AccountNumber: acct,
				SumAmountIn:   tsummary.SumIn,
				SumAmountOut:  tsummary.SumOut,
				SumTotal:      tsummary.SumTotal,
				TransactionDate: &date.Date{
					Year:  int32(now.Year()),
					Month: int32(now.Month()),
					Day:   int32(now.Day()),
				},
			}
			return stream.SendAndClose(&res)
		}
		if err != nil {
			log.Fatalln("Error while reading from client: %v", err)
		}
		acct = req.AccountNumber
		ts, err := toTime(req.GetTimestamp())
		if err != nil {
			log.Fatalln("Error while parsing timestamp %v: %v", req.Timestamp, err)
		}

		ttype := bank.TransactionTypeUnknown
		if req.Type == protobank.TransactionType_TRANSACTION_TYPE_IN {
			ttype = bank.TransactionTypeIn
		} else if req.Type == protobank.TransactionType_TRANSACTION_TYPE_OUT {
			ttype = bank.TransactionTypeOut

		}

		tcur := bank.Transaction{
			Amount:          req.Amount,
			Timestamp:       ts,
			TransactionType: ttype,
		}
		_, err = grpcAdapter.bankService.CreateTransaction(acct, tcur)
		if err != nil {
			log.Fatalf("Error while creating transaction: %v", err)
		}
		err = grpcAdapter.bankService.CalculateTransactionSummary(&tsummary, tcur)
		if err != nil {
			return err
		}
	}

}

func currentDateTime() *datetime.DateTime {

	now := time.Now()
	return &datetime.DateTime{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		Day:        int32(now.Day()),
		Hours:      int32(now.Hour()),
		Minutes:    int32(now.Minute()),
		Seconds:    int32(now.Second()),
		Nanos:      int32(now.Nanosecond()),
		TimeOffset: &datetime.DateTime_UtcOffset{},
	}
}

func (grpcAdapter *GRPCAdapter) TransferMultiple(stream grpc.ClientStreamingServer[protobank.TransferRequest, protobank.TransferResponse]) error {
	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			log.Println("Client cancelled stream")
			return nil
		default:
			req, err := stream.Recv()

			if err == io.EOF {
				return nil
			}
			if err != nil {
				log.Fatalf("Error while reading from client:%v\n", err)
			}

			tt := bank.TransferTransaction{
				FromAccount: req.FromAccountNumber,
				ToAccount:   req.ToAccountNumber,
				Currency:    req.Current,
				Amount:      float64(req.Amount),
			}

			_, transferSuccess, err := grpcAdapter.bankService.Transfer(tt)
			if err != nil {
				return err
			}

			res := protobank.TransferResponse{
				FromAccountNumber: req.FromAccountNumber,
				ToAccountNumber:   req.ToAccountNumber,
				Amount:            float64(req.Amount),
				Timestamp:         currentDateTime(),
				Current:           req.Current,
			}

			if transferSuccess {
				res.Status = protobank.TransferStatus_TRANSFER_STATUS_SUCCESS
			} else {
				res.Status = protobank.TransferStatus_TRANSFER_STATUS_FAILED
			}
			err = stream.SendMsg(&res)
			if err != nil {
				log.Fatalf("Error sending response to client :%v\n", err)
			}
		}
	}
}
