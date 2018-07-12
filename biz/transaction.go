package biz

import (
	"github.com/dongdongking008/softtrans/model"
	"github.com/dongdongking008/softtrans/dao"
	"github.com/cuigh/auxo/log"
)

// Event return a event biz instance.
var Transaction = &transBiz{}

type transBiz struct {
}

func (t *transBiz) Create(trans *model.Transaction) (err error) {
	do(func(d *dao.Dao) {
		err = d.TransCreate(trans)
		if err != nil {
			log.Get("transaction").Errorf("Create trans `%+v` failed: %v", trans, err)
		}
	})
	return err
}

func (t *transBiz) TransGet(transUniqId string) (trans *model.Transaction, err error) {
	do(func(d *dao.Dao) {
		trans, err = d.TransGet(transUniqId)
	})
	return
}

func (t *transBiz) TransGetByTransId(transId *model.TransactionId) (trans *model.Transaction, err error) {
	do(func(d *dao.Dao) {
		trans, err = d.TransGetByTransId(transId)
	})
	return
}

func (t *transBiz) AddStep(transUniqId string, step *model.TransactionStep) (err error) {
	do(func(d *dao.Dao) {
		err = d.TransStepAdd(transUniqId, step)
	})
	return
}

func (t *transBiz) TransConfirm(transUniqId string) (err error) {
	do(func(d *dao.Dao) {
		err = d.TransConfirm(transUniqId)
	})
	return
}

func (t *transBiz) TransCancel(transUniqId string) (err error) {
	do(func(d *dao.Dao) {
		err = d.TransCancel(transUniqId)
	})
	return
}

func (t *transBiz) TransCancelSuccess(transUniqId string) (err error) {
	do(func(d *dao.Dao) {
		err = d.TransCancelSuccess(transUniqId)
	})
	return
}

func (t *transBiz) TransGetExpiredList(topN int32) (transUniqIds []string, err error) {
	do(func(d *dao.Dao) {
		transUniqIds, err = d.TransGetExpiredList(topN)
	})
	return
}

func (t *transBiz) TransGetRollingBackList(topN int32) (transList []*model.Transaction, err error) {
	do(func(d *dao.Dao) {
		transList, err = d.TransGetRollingBackList(topN)
	})
	return
}