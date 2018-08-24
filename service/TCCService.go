package service

import (
	"context"
	"github.com/cuigh/auxo/errors"
	"github.com/dongdongking008/softtrans/biz"
	"github.com/dongdongking008/softtrans/contract"
	"github.com/dongdongking008/softtrans/model"
	"github.com/dongdongking008/softtrans/util/clientname"
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
		Status: model.TransactionStatusTry,
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
			return &contract.BeginTransResponse{TransUniqId: trans.ID.Hex()}, nil
		}
	}

	return nil, err

}

func (s *TCCService) TryStep(ctx context.Context, request *contract.TryStepRequest) (*contract.TryStepResponse, error) {
	if request.GetTransUniqId() == "" {
		return nil, errors.Coded(int32(contract.TryStepResponse_EmptyTransUniqId), "TransUniqId is Empty!")
	}
	step := request.GetStep()
	if step == nil || step.GetStepId() == "" || step.GetServerName() == "" ||
		step.GetServiceName() == "" || step.GetConfirmMethodName() == "" ||
		step.GetCancelMethodName() == "" {
		return nil, errors.Coded(int32(contract.TryStepResponse_InvalidStepInfo), "Step is invalid!")
	}
	transStep := &model.TransactionStep{
		StepId:            step.GetStepId(),
		Args:              step.GetArgs(),
		ServerName:        step.GetServerName(),
		ServiceName:       step.GetServiceName(),
		ConfirmMethodName: step.GetConfirmMethodName(),
		CancelMethodName:  step.GetCancelMethodName(),
		ClientName:        step.GetClientName(),
	}
	err := biz.Transaction.AddStep(request.GetTransUniqId(), transStep)
	if err == nil {
		return &contract.TryStepResponse{TransUniqId: request.GetTransUniqId(), StepId: step.GetStepId()}, nil
	}
	return nil, err
}

func (s *TCCService) ConfirmTrans(ctx context.Context, request *contract.ConfirmTransRequest) (*contract.ConfirmTransResponse, error) {
	if request.GetTransUniqId() == "" {
		return nil, errors.Coded(int32(contract.ConfirmTransResponse_EmptyTransUniqId), "TransUniqId is Empty!")
	}
	err := biz.Transaction.TransConfirm(request.GetTransUniqId())
	if err == nil {
		return &contract.ConfirmTransResponse{TransUniqId: request.GetTransUniqId()}, nil
	}
	return nil, err
}

func (s *TCCService) ConfirmTransSuccess(ctx context.Context, request *contract.ConfirmTransRequest) (*contract.ConfirmTransSuccessResponse, error) {
	if request.GetTransUniqId() == "" {
		return nil, errors.Coded(int32(contract.ConfirmTransSuccessResponse_EmptyTransUniqId), "TransUniqId is Empty!")
	}
	err := biz.Transaction.TransCancel(request.GetTransUniqId())
	if err == nil {
		return &contract.ConfirmTransSuccessResponse{TransUniqId: request.GetTransUniqId()}, nil
	}
	return nil, err
}

func (s *TCCService) CancelTrans(ctx context.Context, request *contract.CancelTransRequest) (*contract.CancelTransResponse, error) {
	if request.GetTransUniqId() == "" {
		return nil, errors.Coded(int32(contract.CancelTransResponse_EmptyTransUniqId), "TransUniqId is Empty!")
	}
	err := biz.Transaction.TransCancel(request.GetTransUniqId())
	if err == nil {
		return &contract.CancelTransResponse{TransUniqId: request.GetTransUniqId()}, nil
	}
	return nil, err
}

func (s *TCCService) CancelTransSuccess(ctx context.Context, request *contract.CancelTransRequest) (*contract.CancelTransSuccessResponse, error) {
	if request.GetTransUniqId() == "" {
		return nil, errors.Coded(int32(contract.CancelTransSuccessResponse_EmptyTransUniqId), "TransUniqId is Empty!")
	}
	err := biz.Transaction.TransCancel(request.GetTransUniqId())
	if err == nil {
		return &contract.CancelTransSuccessResponse{TransUniqId: request.GetTransUniqId()}, nil
	}
	return nil, err
}

func (s *TCCService) GetTrans(ctx context.Context, request *contract.GetTransRequest) (*contract.GetTransResponse, error) {
	transUniqId := request.GetTransUniqId()
	if transUniqId == "" {
		return nil, errors.Coded(int32(contract.GetTransResponse_EmptyTransUniqId), "TransUniqId is Empty!")
	}
	trans, err := biz.Transaction.TransGet(transUniqId)
	if err == nil {
		return &contract.GetTransResponse{Transaction: transModelToProto(trans)}, nil
	} else {
		return nil, err
	}
}

// Get expired transactions
func (s *TCCService) GetExpiredTransList(ctx context.Context, request *contract.GetExpiredTransListRequest) (*contract.GetExpiredTransListResponse, error) {
	transUniqIds, err := biz.Transaction.TransGetExpiredList(request.GetTopN())
	if err == nil {
		return &contract.GetExpiredTransListResponse{TransUniqIds: transUniqIds}, nil
	} else {
		return nil, err
	}
}

// Get confirming transactions
func (s *TCCService) GetConfirmingList(ctx context.Context, request *contract.GetConfirmingTransListRequest) (*contract.GetConfirmingTransListResponse, error) {
	transList, err := biz.Transaction.TransGetConfirmingList(request.GetTopN())
	if err == nil {
		return &contract.GetConfirmingTransListResponse{Transactions: transModelsToProtos(transList)}, nil
	} else {
		return nil, err
	}
}

// Get cancelling transactions
func (s *TCCService) GetCancellingTransList(ctx context.Context, request *contract.GetCancellingTransListRequest) (*contract.GetCancellingTransListResponse, error) {
	transList, err := biz.Transaction.TransGetCancellingList(request.GetTopN())
	if err == nil {
		return &contract.GetCancellingTransListResponse{Transactions: transModelsToProtos(transList)}, nil
	} else {
		return nil, err
	}
}

func transModelsToProtos(transList []*model.Transaction) []*contract.Transaction {
	if transList == nil {
		return nil
	}
	transactions := make([]*contract.Transaction, 0, 5)
	for _, trans := range transList {
		transactions = append(transactions, transModelToProto(trans))
	}
	return transactions
}

func transModelToProto(trans *model.Transaction) *contract.Transaction {

	if trans == nil {
		return nil
	}

	steps := make([]*contract.TransactionStep, 0, 5)
	for _, step := range trans.Steps {
		steps = append(steps, &contract.TransactionStep{
			StepId:            step.StepId,
			Args:              step.Args,
			ServerName:        step.ServerName,
			ServiceName:       step.ServiceName,
			ConfirmMethodName: step.ConfirmMethodName,
			CancelMethodName:  step.CancelMethodName,
			ClientName:        step.ClientName,
		})
	}

	return &contract.Transaction{
		TransactionId: &contract.TransactionId{
			AppId:   trans.TransId.AppId,
			BusCode: trans.TransId.BusCode,
			TrxId:   trans.TransId.TrxId,
		},
		TransUniqId: trans.ID.Hex(),
		Steps:       steps,
		Status:      contract.Transaction_TransactionStatus(int32(trans.Status)),
	}
}
