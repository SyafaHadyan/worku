// Package payment handles the payment processing and interacts with third party API
package payment

import (
	"errors"
	"fmt"
	"log"

	"github.com/SyafaHadyan/worku/internal/domain/dto"
	"github.com/SyafaHadyan/worku/internal/infra/env"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/example"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentItf interface {
	CreateSnapPayment(createMidtransSnapOrder dto.CreateMidtransSnapOrder) (*snap.Response, error)
	CreateCoreAPIPayment(createMidtransCoreAPIOrder dto.CreateMidtransCoreAPIOrder) (*coreapi.ChargeResponse, error)
	CreateCoreAPISubscription(createMidtransCoreAPISubscription dto.CreateMidtransCoreAPISubscription) (*coreapi.CreateSubscriptionResponse, error)
}

type Payment struct {
	Client                      *snap.Client
	CoreAPIClient               *coreapi.Client
	basePrice                   BasePrice
	baseMonthDiscountPercentage BaseMonthDiscountPercentage
}

type BasePrice struct {
	Day   int64
	Month int64
}

type BaseMonthDiscountPercentage struct {
	Month6  int
	Month12 int
}

func New(env *env.Env) *Payment {
	var SnapClient snap.Client
	SnapClient.New(env.MidtransServerKey, midtrans.Sandbox)

	CoreAPIClient := coreapi.Client{
		ServerKey: env.MidtransServerKey,
		ClientKey: env.MidtransClientKey,
		Env:       midtrans.Sandbox,
	}

	BasePrice := BasePrice{
		Day:   1000,
		Month: 30000,
	}

	BaseMonthDiscountPercentage := BaseMonthDiscountPercentage{
		Month6:  42,
		Month12: 50,
	}

	Payment := Payment{
		Client:                      &SnapClient,
		CoreAPIClient:               &CoreAPIClient,
		basePrice:                   BasePrice,
		baseMonthDiscountPercentage: BaseMonthDiscountPercentage,
	}

	return &Payment
}

func (p *Payment) calculatePrice(interval int) int64 {
	result := p.basePrice.Month

	log.Println(interval)

	switch interval {
	case 180:
		result = (result * int64(100-p.baseMonthDiscountPercentage.Month6) / 100) * 6
	case 360:
		result = (result * int64(100-p.baseMonthDiscountPercentage.Month12) / 100) * 12
	default:
		result = p.basePrice.Day * int64(interval)
	}

	return result
}

func (p *Payment) chargeTransactionWithMap() (coreapi.ResponseWithMap, error) {
	res, err := p.CoreAPIClient.ChargeTransactionWithMap(example.CoreParam())
	if err != nil {
		return nil,
			errors.New(fmt.Sprintf("Error coreapi api %s", err.GetMessage()))
	}

	return res, nil
}

func (p *Payment) getCardToken(card dto.Card) (*coreapi.CardTokenResponse, error) {
	midtrans.ClientKey = example.SandboxClientKey2
	res, err := p.CoreAPIClient.CardToken(
		card.CardNumber,
		card.ExpiryMonth,
		card.ExpiryYear,
		card.CVV,
		p.CoreAPIClient.ClientKey,
	)
	if err != nil {
		return nil,
			errors.New(fmt.Sprintf("Error get card token %s", err.GetMessage()))
	}

	return res, nil
}

func (p *Payment) registerCard(card dto.Card) (*coreapi.CardRegisterResponse, error) {
	res, err := p.CoreAPIClient.RegisterCard(
		card.CardNumber,
		card.ExpiryMonth,
		card.ExpiryYear,
		p.CoreAPIClient.ClientKey,
	)
	if err != nil {
		return &coreapi.CardRegisterResponse{},
			errors.New(fmt.Sprintf("Error register card token %s", err.GetMessage()))
	}

	return res, err
}

func (p *Payment) cardPointInquiry(card dto.Card) (*coreapi.PointInquiryResponse, error) {
	cardToken, err := p.getCardToken(card)
	if err != nil {
		return nil, err
	}

	res, midtransErr := p.CoreAPIClient.CardPointInquiry(cardToken.TokenID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error card point inquiry %s", midtransErr.GetMessage()))
	}

	return res, err
}

func (p *Payment) getBin(bin string) (*coreapi.BinResponse, error) {
	res, err := p.CoreAPIClient.GetBIN(bin)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error get bin %s", err.GetMessage()))
	}

	return res, nil
}

func (p *Payment) requestCreditCard(createMidtransCoreAPIOrder dto.CreateMidtransCoreAPIOrder) (*coreapi.ChargeResponse, error) {
	grossAmount := p.calculatePrice(createMidtransCoreAPIOrder.Interval)
	tokenID, err := p.getCardToken(createMidtransCoreAPIOrder.Card)
	if err != nil {
		return nil, err
	}

	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeCreditCard,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  createMidtransCoreAPIOrder.TransactionDetails.OrderID,
			GrossAmt: grossAmount,
		},
		CreditCard: &coreapi.CreditCardDetails{
			TokenID:        tokenID.TokenID,
			Authentication: true,
		},
	}

	res, _ := p.CoreAPIClient.ChargeTransaction(chargeReq)

	return res, nil
}

func (p *Payment) CreateSnapPayment(createMidtransSnapOrder dto.CreateMidtransSnapOrder) (*snap.Response, error) {
	grossAmount := p.calculatePrice(createMidtransSnapOrder.Interval)

	snapRequest := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  createMidtransSnapOrder.TransactionDetails.OrderID,
			GrossAmt: grossAmount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: createMidtransSnapOrder.CustomerDetail.FirstName,
			Email: createMidtransSnapOrder.CustomerDetail.Email,
		},
	}

	res, err := p.Client.CreateTransaction(snapRequest)

	return res, err
}

func (p *Payment) CreateCoreAPIPayment(createMidtransCoreAPIOrder dto.CreateMidtransCoreAPIOrder) (*coreapi.ChargeResponse, error) {
	var res *coreapi.ChargeResponse
	var err error

	if createMidtransCoreAPIOrder.PaymentMethod == "CREDIT_CARD" {
		res, err = p.requestCreditCard(createMidtransCoreAPIOrder)
		log.Println(res)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return res, nil
}

func (p *Payment) CreateCoreAPISubscription(createMidtransCoreAPISubscription dto.CreateMidtransCoreAPISubscription) (*coreapi.CreateSubscriptionResponse, error) {
	var amount int64
	subscriptionName := fmt.Sprintf(
		"WorkU-%d-%s",
		createMidtransCoreAPISubscription.Interval,
		createMidtransCoreAPISubscription.IntervalUnit,
	)

	if createMidtransCoreAPISubscription.IntervalUnit == "month" {
		amount = p.calculatePrice(createMidtransCoreAPISubscription.Interval)
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
		Currency:    createMidtransCoreAPISubscription.Currency,
		PaymentType: coreapi.PaymentTypeCreditCard,
		Token:       "TODO",
		Schedule: coreapi.ScheduleDetails{
			Interval:     createMidtransCoreAPISubscription.Interval,
			IntervalUnit: createMidtransCoreAPISubscription.IntervalUnit,
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: createMidtransCoreAPISubscription.CustomerDetail.FirstName,
			Email: createMidtransCoreAPISubscription.CustomerDetail.Email,
		},
	}

	res, err := p.CoreAPIClient.CreateSubscription(req)

	return res, err
}
