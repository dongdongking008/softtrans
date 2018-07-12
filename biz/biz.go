package biz

import (
	"github.com/dongdongking008/softtrans/dao"
	"github.com/cuigh/auxo/errors"
)

func do(fn func(d *dao.Dao)) {
	d, err := dao.Get()
	if err != nil {
		panic(errors.Wrap(err, "failed to load storage engine"))
	}

	fn(d)
}