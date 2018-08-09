package model

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type TransactionStatus int32

const (
	TransactionStatusTry        TransactionStatus = 0
	TransactionStatusConfirming TransactionStatus = 10
	TransactionStatusConfirmed  TransactionStatus = 20
	TransactionStatusCancelling TransactionStatus = 30
	TransactionStatusCancelled  TransactionStatus = 40
)

type Transaction struct {
	ID             bson.ObjectId     `bson:"_id"`
	TransId        TransactionId     `bson:"trans_id" json:"trans_id,omitempty"`
	Steps          []TransactionStep `bson:"steps" json:"steps,omitempty"`
	Status         TransactionStatus `bson:"status" json:"status,omitempty"`
	EnterTime      time.Time         `bson:"enter_time"`
	ExpireTime     time.Time         `bson:"expire_time"`
	LastUpdateTime time.Time         `bson:"lu_time"`
}

type TransactionId struct {
	AppId   string `bson:"app_id" json:"app_id,omitempty"`
	BusCode string `bson:"bus_code" json:"bus_code,omitempty"`
	TrxId   string `bson:"trx_id" json:"trx_id,omitempty"`
}

type TransactionStep struct {
	StepId            string `bson:"step_id" json:"step_id,omitempty"`
	Args              []byte `bson:"args" json:"args,omitempty"`
	ServerName        string `bson:"server_name" json:"server_name,omitempty"`
	ServiceName       string `bson:"service_name" json:"service_name,omitempty"`
	ConfirmMethodName string `bson:"confirm_method_name" json:"confirm_method_name,omitempty"`
	CancelMethodName  string `bson:"cancel_method_name" json:"cancel_method_name,omitempty"`
}
