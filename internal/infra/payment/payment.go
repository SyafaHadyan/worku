// Package payment handles the payment processing and interacts with third party API
package payment

import (
	"fmt"

	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentItf interface {
	CreatePayment(createMidtransOrder dto.CreateMidtransOrder) (*snap.Response, error)
	CreateSubscription(createMidtransSubscription dto.CreateMidtransSubscription) (*coreapi.CreateSubscriptionResponse, error)
}

type Payment struct {
	Client                      *snap.Client
	CoreAPIClient               *coreapi.Client
	BasePrice                   BasePrice
	BaseMonthDiscountPercentage BaseMonthDiscountPercentage
}

type BasePrice struct {
	Month int64
}

type BaseMonthDiscountPercentage struct {
	Month6  int
	Month12 int
}

func New(env *env.Env) *Payment {
	var client snap.Client
	client.New(env.MidtransServerKey, midtrans.Sandbox)

	BasePrice := BasePrice{
		Month: 30000,
	}

	BaseMonthDiscountPercentage := BaseMonthDiscountPercentage{
		Month6:  42,
		Month12: 50,
	}

	Payment := Payment{
		Client:                      &client,
		CoreAPIClient:               &coreapi.Client{},
		BasePrice:                   BasePrice,
		BaseMonthDiscountPercentage: BaseMonthDiscountPercentage,
	}

	return &Payment
}

func (p *Payment) calculateSubscriptionPrice(interval int) int64 {
	result := p.BasePrice.Month

	switch interval {
	case 180:
		result = (result * int64(100-p.BaseMonthDiscountPercentage.Month6) / 100) * 6
	case 360:
		result = (result * int64(100-p.BaseMonthDiscountPercentage.Month12) / 100) * 12
	}

	return result
}

func (p *Payment) CreatePayment(createMidtransOrder dto.CreateMidtransOrder) (*snap.Response, error) {
	grossAmount := p.calculateSubscriptionPrice(createMidtransOrder.Interval)

	snapRequest := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  createMidtransOrder.TransactionDetails.OrderID,
			GrossAmt: grossAmount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: createMidtransOrder.CustomerDetail.FirstName,
			Email: createMidtransOrder.CustomerDetail.Email,
		},
	}

	res, err := p.Client.CreateTransaction(snapRequest)

	return res, err
}

func (p *Payment) CreateSubscription(createMidtransSubscription dto.CreateMidtransSubscription) (*coreapi.CreateSubscriptionResponse, error) {
	var amount int64
	subscriptionName := fmt.Sprintf(
		"WorkU-%d-%s",
		createMidtransSubscription.Interval,
		createMidtransSubscription.IntervalUnit,
	)

	if createMidtransSubscription.IntervalUnit == "month" {
		amount = p.calculateSubscriptionPrice(createMidtransSubscription.Interval)
	} else {
		/*
		* For testing only
		* Duration: 1 day
		* Amount: Rp1
		 */
		amount = 1
	}

	// TODO: retrieve token from user's saved payment method (credit card only)
	req := &coreapi.SubscriptionReq{
		Name:        subscriptionName,
		Amount:      amount,
		Currency:    createMidtransSubscription.Currency,
		PaymentType: coreapi.PaymentTypeCreditCard,
		Token:       "todo",
		Schedule: coreapi.ScheduleDetails{
			Interval:     createMidtransSubscription.Interval,
			IntervalUnit: createMidtransSubscription.IntervalUnit,
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: createMidtransSubscription.CustomerDetail.FirstName,
			Email: createMidtransSubscription.CustomerDetail.Email,
		},
	}

	res, err := p.CoreAPIClient.CreateSubscription(req)

	return res, err
}
