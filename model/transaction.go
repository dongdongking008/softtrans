package model

import (
	"time"
	"github.com/globalsign/mgo/bson"
)

type TransactionStatus int32

const (
	TransactionStatusInit TransactionStatus = 0
	TransactionStatusCommitted TransactionStatus = 10
	TransactionStatusRollingBack TransactionStatus = 20
	TransactionStatusRolledBack TransactionStatus = 40
)

type Transaction struct {
	ID       bson.ObjectId      `bson:"_id"`
	TransId TransactionId   `bson:"trans_id" json:"trans_id,omitempty"`
	Steps []TransactionStep   `bson:"steps" json:"steps,omitempty"`
	Status TransactionStatus `bson:"status" json:"status,omitempty"`
	EnterTime     time.Time   `bson:"enter_time"`
	ExpireTime	  time.Time   `bson:"expire_time"`
	LastUpdateTime time.Time	`bson:"lu_time"`
}

type TransactionId struct {
	AppId string   `bson:"app_id" json:"app_id,omitempty"`
	BusCode string   `bson:"bus_code" json:"bus_code,omitempty"`
	TrxId string   `bson:"trx_id" json:"trx_id,omitempty"`
}

type TransactionStep struct {
	StepId string   `bson:"step_id" json:"step_id,omitempty"`
	Args []byte   `bson:"args" json:"args,omitempty"`
}
