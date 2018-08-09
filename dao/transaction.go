package dao

import (
	"fmt"
	"github.com/cuigh/auxo/db/mongo"
	"github.com/cuigh/auxo/errors"
	"github.com/dongdongking008/softtrans/contract"
	"github.com/dongdongking008/softtrans/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"time"
)

func (d *Dao) TransGet(transUniqId string) (trans *model.Transaction, err error) {
	err = d.do(func(db mongo.DB) error {
		trans = &model.Transaction{}

		q := bson.M{
			"_id": transUniqId,
		}
		errDB := db.C("transaction").Find(q).One(trans)
		if errDB == mgo.ErrNotFound {
			trans = nil
			return errDB
		} else if err != nil {
			trans = nil
		}
		return errDB
	})
	return
}

func (d *Dao) TransGetByTransId(id *model.TransactionId) (trans *model.Transaction, err error) {
	err = d.do(func(db mongo.DB) error {
		trans = &model.Transaction{}

		q := bson.M{
			"trans_id.app_id":   id.AppId,
			"trans_id.bus_code": id.BusCode,
			"trans_id.trx_id":   id.TrxId,
		}
		errDB := db.C("transaction").Find(q).One(trans)
		if errDB == mgo.ErrNotFound {
			trans = nil
			return nil
		} else if errDB != nil {
			trans = nil
		}
		return errDB
	})
	return
}

func (d *Dao) TransCreate(trans *model.Transaction) (err error) {
	err = d.do(func(db mongo.DB) error {
		return db.C("transaction").Insert(trans)
	})

	if err != nil && mgo.IsDup(err) {
		err = errors.Coded(int32(contract.BeginTransResponse_DuplicateRequest), err.Error())
	}

	return
}

func (d *Dao) TransStepAdd(transUniqId string, step *model.TransactionStep) (err error) {
	err = d.do(func(db mongo.DB) error {
		return db.C("transaction").Update(bson.D{
			{Name: "_id", Value: transUniqId},
			{Name: "status", Value: model.TransactionStatusInit},
		},
			bson.D{
				{Name: "$addToSet", Value: bson.D{{Name: "steps", Value: step}}},
				{Name: "$set", Value: bson.D{{Name: "lu_time", Value: time.Now()}}},
			})
	})
	if err == mgo.ErrNotFound {
		trans, errDB := d.TransGet(transUniqId)
		if trans != nil {
			err = errors.Coded(int32(contract.TryStepResponse_TransactionStatusError), fmt.Sprintf("Transaction status is %d", trans.Status))
		} else if errDB == mgo.ErrNotFound {
			err = errors.Coded(int32(contract.TryStepResponse_TransactionNotFound), err.Error())
		} else {
			err = errDB
		}
	}
	return
}

func (d *Dao) TransConfirm(transUniqId string) (err error) {
	err = d.do(func(db mongo.DB) error {
		return db.C("transaction").Update(bson.D{
			{Name: "_id", Value: transUniqId},
			{Name: "status", Value: model.TransactionStatusInit},
		},
			bson.D{
				{Name: "status", Value: model.TransactionStatusCommitted},
				{Name: "$set", Value: bson.D{{Name: "lu_time", Value: time.Now()}}},
			})
	})
	if err == mgo.ErrNotFound {
		trans, errDB := d.TransGet(transUniqId)
		if trans != nil {
			if trans.Status == model.TransactionStatusCommitted {
				err = nil
			} else {
				err = errors.Coded(int32(contract.ConfirmTransResponse_TransactionStatusError), fmt.Sprintf("Transaction status is %d", trans.Status))
			}
		} else if errDB == mgo.ErrNotFound {
			err = errors.Coded(int32(contract.ConfirmTransResponse_TransactionNotFound), err.Error())
		} else {
			err = errDB
		}
	}
	return
}

func (d *Dao) TransCancel(transUniqId string) (err error) {
	err = d.do(func(db mongo.DB) error {
		return db.C("transaction").Update(bson.D{
			{Name: "_id", Value: transUniqId},
			{Name: "status", Value: model.TransactionStatusInit},
		},
			bson.D{
				{Name: "status", Value: model.TransactionStatusRollingBack},
				{Name: "$set", Value: bson.D{{Name: "lu_time", Value: time.Now()}}},
			})
	})
	if err == mgo.ErrNotFound {
		trans, errDB := d.TransGet(transUniqId)
		if trans != nil {
			if trans.Status == model.TransactionStatusRollingBack ||
				trans.Status == model.TransactionStatusRolledBack {
				err = nil
			} else {
				err = errors.Coded(int32(contract.CancelTransResponse_TransactionStatusError), fmt.Sprintf("Transaction status is %d", trans.Status))
			}
		} else if errDB == mgo.ErrNotFound {
			err = errors.Coded(int32(contract.CancelTransResponse_TransactionNotFound), err.Error())
		} else {
			err = errDB
		}
	}
	return
}

func (d *Dao) TransCancelSuccess(transUniqId string) (err error) {
	err = d.do(func(db mongo.DB) error {
		return db.C("transaction").Update(bson.D{
			{Name: "_id", Value: transUniqId},
			{Name: "status", Value: model.TransactionStatusRollingBack},
		},
			bson.D{
				{Name: "status", Value: model.TransactionStatusRolledBack},
				{Name: "$set", Value: bson.D{{Name: "lu_time", Value: time.Now()}}},
			})
	})
	if err == mgo.ErrNotFound {
		trans, errDB := d.TransGet(transUniqId)
		if trans != nil {
			if trans.Status == model.TransactionStatusRolledBack {
				err = nil
			} else {
				err = errors.Coded(int32(contract.CancelTransSuccessResponse_TransactionStatusError), fmt.Sprintf("Transaction status is %d", trans.Status))
			}
		} else if errDB == mgo.ErrNotFound {
			err = errors.Coded(int32(contract.CancelTransSuccessResponse_TransactionNotFound), err.Error())
		} else {
			err = errDB
		}
	}
	return
}

func (d *Dao) TransGetExpiredList(topN int32) (transUniqIds []string, err error) {
	err = d.do(func(db mongo.DB) error {
		transUniqIds = []string{}
		transList := make([]*model.Transaction, 0)
		q := bson.M{
			"expire_time": bson.M{"$lt": time.Now()},
			"status":      model.TransactionStatusInit,
		}
		s := bson.M{"_id": 1}
		errDB := db.C("transaction").Find(q).Select(s).
			Sort("_id").Limit(int(topN)).All(&transList)
		if errDB == nil {
			for _, trans := range transList {
				transUniqIds = append(transUniqIds, trans.ID.String())
			}
		}
		return errDB
	})
	return
}

func (d *Dao) TransGetRollingBackList(topN int32) (transList []*model.Transaction, err error) {
	err = d.do(func(db mongo.DB) error {
		transList = []*model.Transaction{}
		q := bson.M{
			"status": model.TransactionStatusRollingBack,
		}
		return db.C("transaction").Find(q).Sort("_id").
			Limit(int(topN)).All(&transList)
	})
	return
}
