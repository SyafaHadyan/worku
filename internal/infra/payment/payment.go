// Package payment handles the payment processing and interacts with third party API
package payment

import (
	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentItf interface {
	GenerateSnapRequest(createMidtransOrder dto.CreateMidtransOrder) *snap.Request
	CreatePayment(createMidtransOrder dto.CreateMidtransOrder) (*snap.Response, error)
}

type Payment struct {
	Client *snap.Client
}

func New(env *env.Env) *Payment {
	var client snap.Client
	client.New(env.MidtransServerKey, midtrans.Sandbox)

	Payment := Payment{
		Client: &client,
	}

	return &Payment
}

func (p *Payment) GenerateSnapRequest(createMidtransOrder dto.CreateMidtransOrder) *snap.Request {
	return &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  createMidtransOrder.TransactionDetails.OrderID,
			GrossAmt: int64(createMidtransOrder.TransactionDetails.GrossAmount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: createMidtransOrder.CustomerDetail.FirstName,
			Email: createMidtransOrder.CustomerDetail.Email,
		},
	}
}

func (p *Payment) CreatePayment(createMidtransOrder dto.CreateMidtransOrder) (*snap.Response, error) {
	snapRequest := p.GenerateSnapRequest(createMidtransOrder)

	res, err := p.Client.CreateTransaction(snapRequest)

	return res, err
}
