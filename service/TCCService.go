package service

import (
	"context"
	"github.com/dongdongking008/softtrans/contract"
	"github.com/cuigh/auxo/errors"
	"github.com/dongdongking008/softtrans/util/clientname"
	"github.com/dongdongking008/softtrans/biz"
	"github.com/dongdongking008/softtrans/model"
	"time"
)

const (
	tccTransactionDefaultExpireTime = time.Duration(30)
)

type TCCService struct {

}

// Begin a transaction
func (s *TCCService) BeginTrans(ctx context.Context, request *contract.BeginTransRequest) (*contract.BeginTransResponse, error) {

	transactionId := request.GetTransactionId()
	if transactionId.GetAppId() == "" || clientname.FromContext(ctx) != transactionId.GetAppId() {
		return nil, errors.Coded(int32(contract.BeginTransResponse_InvalidAppId), "Invalid AppId!")
	}

	if transactionId.GetBusCode() == "" {
		return nil, errors.Coded(int32(contract.BeginTransResponse_EmptyBusCode), "BusCode is Empty!")
	}

	if request.GetFailFast() && transactionId.GetTrxId() == "" {
		return nil, errors.Coded(int32(contract.BeginTransResponse_EmptyTRXId), "TrxId is Empty!")
	}

	trans := &model.Transaction{
		Status: model.TransactionStatusInit,
	}
	trans.TransId.AppId = transactionId.GetAppId()
	trans.TransId.BusCode = transactionId.GetBusCode()
	trans.TransId.TrxId = transactionId.GetTrxId()

	trans.EnterTime = time.Now()
	trans.LastUpdateTime = trans.EnterTime

	if request.ExpireTimeSeconds > 0 {
		trans.ExpireTime = trans.EnterTime.Add(time.Second * time.Duration(request.GetExpireTimeSeconds()))
	} else {
		trans.ExpireTime = trans.EnterTime.Add(time.Second * tccTransactionDefaultExpireTime)
	}

	err := biz.Transaction.Create(trans)
	if err != nil {
		if errCoded, ok := err.(*errors.CodedError); ok {
			if errCoded.Code == int32(contract.BeginTransResponse_DuplicateRequest) {
				err = nil
			}
		}
	}

	if err == nil {
		trans, err = biz.Transaction.TransGetByTransId(&trans.TransId)
		if err == nil {
			return &contract.BeginTransResponse{ TransUniqId: trans.ID.String() }, nil
		}
	}

	return nil, err

}

func (s *TCCService) TryStep(ctx context.Context, request *contract.TryStepRequest) (*contract.TryStepResponse, error) {
	if request.GetTransUniqId() == "" {
		return nil, errors.Coded(int32(contract.TryStepResponse_EmptyTransUniqId), "TransUniqId is Empty!")
	}
	if request.GetStepId() == "" {
		return nil, errors.Coded(int32(contract.TryStepResponse_EmptyStepId), "StepId is Empty!")
	}
	transStep := &model.TransactionStep{
		StepId: request.GetStepId(),
		Args: request.GetArgs(),
	}
	err := biz.Transaction.AddStep(request.GetTransUniqId(), transStep)
	if err == nil {
		return &contract.TryStepResponse{ TransUniqId:request.GetTransUniqId(), StepId: request.GetStepId()}, nil
	}
	return nil, err
}

func (s *TCCService) ConfirmTrans(ctx context.Context, request *contract.ConfirmTransRequest) (*contract.ConfirmTransResponse, error) {
	if request.GetTransUniqId() == "" {
		return nil, errors.Coded(int32(contract.ConfirmTransResponse_EmptyTransUniqId), "TransUniqId is Empty!")
	}
	err := biz.Transaction.TransConfirm(request.GetTransUniqId())
	if err == nil {
		return &contract.ConfirmTransResponse{ TransUniqId:request.GetTransUniqId()}, nil
	}
	return nil, err
}

func (s *TCCService) CancelTrans(ctx context.Context, request *contract.CancelTransRequest) (*contract.CancelTransResponse, error) {
	if request.GetTransUniqId() == "" {
		return nil, errors.Coded(int32(contract.CancelTransResponse_EmptyTransUniqId), "TransUniqId is Empty!")
	}
	err := biz.Transaction.TransCancel(request.GetTransUniqId())
	if err == nil {
		return &contract.CancelTransResponse{ TransUniqId:request.GetTransUniqId()}, nil
	}
	return nil, err
}

func (s *TCCService) CancelTransSuccess(ctx context.Context, request *contract.CancelTransRequest) (*contract.CancelTransSuccessResponse, error) {
	if request.GetTransUniqId() == "" {
		return nil, errors.Coded(int32(contract.CancelTransSuccessResponse_EmptyTransUniqId), "TransUniqId is Empty!")
	}
	err := biz.Transaction.TransCancel(request.GetTransUniqId())
	if err == nil {
		return &contract.CancelTransSuccessResponse{ TransUniqId:request.GetTransUniqId()}, nil
	}
	return nil, err
}

func (s *TCCService) GetTrans(ctx context.Context, request *contract.GetTransRequest) (*contract.GetTransResponse, error) {
	transactionId := request.GetTransUniqId()
	trans, err := biz.Transaction.TransGet(transactionId)
	if err == nil {
		return &contract.GetTransResponse{ Transaction: transModelToProto(trans) }, nil
	} else {
		return nil, err
	}
}

// Get expired transactions
func (s *TCCService) GetExpiredTransList(ctx context.Context, request *contract.GetExpiredTransListRequest) (*contract.GetExpiredTransListResponse, error) {
	transUniqIds, err := biz.Transaction.TransGetExpiredList(request.GetTopN())
	if err == nil {
		return &contract.GetExpiredTransListResponse{ TransUniqIds: transUniqIds }, nil
	} else {
		return nil, err
	}
}
// Get rolling back transactions
func (s *TCCService) GetRollingBackTransList(ctx context.Context, request *contract.GetRollingBackTransListRequest) (*contract.GetRollingBackTransListResponse, error) {
	transList, err := biz.Transaction.TransGetRollingBackList(request.GetTopN())
	if err == nil {
		return &contract.GetRollingBackTransListResponse{ Transactions: transModelsToProtos(transList) }, nil
	} else {
		return nil, err
	}
}

func transModelsToProtos(transList []*model.Transaction) []*contract.Transaction {
	if transList == nil {
		return nil
	}
	transactions := make([]*contract.Transaction, 5)
	for _, trans := range transList {
		transactions = append(transactions, transModelToProto(trans))
	}
	return transactions
}

func transModelToProto(trans *model.Transaction) *contract.Transaction {

	if trans == nil {
		return nil
	}

	steps := make([]*contract.TransactionStep, 5)
	for _, step := range trans.Steps  {
		steps = append(steps, &contract.TransactionStep{
			StepId: step.StepId,
			Args: step.Args,
		})
	}

	return &contract.Transaction{
		TransactionId: &contract.TransactionId {
			AppId: trans.TransId.AppId,
			BusCode: trans.TransId.BusCode,
			TrxId: trans.TransId.TrxId,
		},
		TransUniqId: trans.ID.String(),
		Steps: steps,
		Status: contract.Transaction_TransactionStatus(int32(trans.Status)),
	}
}
