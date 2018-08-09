package dao

import (
	"github.com/cuigh/auxo/db/mongo"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/util/lazy"
	"github.com/globalsign/mgo"
)

const (
	DBName  = "tcc"
	PKGNAME = "softtrans.dao"
)

var (
	indexes = map[string][]mgo.Index{
		"transaction": {
			mgo.Index{Key: []string{"trans_id"}, Unique: true},
			mgo.Index{Key: []string{"status"}},
			mgo.Index{Key: []string{"enter_time"}},
			mgo.Index{Key: []string{"expire_time"}},
			mgo.Index{Key: []string{"lu_time"}},
		},
	}
)

var (
	value = lazy.Value{New: create}
)

func Get() (*Dao, error) {
	v, err := value.Get()
	if err != nil {
		return nil, err
	}
	return v.(*Dao), nil
}

type Dao struct {
	dbName string
	logger log.Logger
}

func (d *Dao) Init() error {
	return d.do(func(db mongo.DB) error {
		for name, ins := range indexes {
			c := db.C(name)
			for _, in := range ins {
				err := c.EnsureIndex(in)
				if err != nil {
					d.logger.Warnf("Ensure index %s-%v failed: %v", name, in.Key, err)
				}
			}
		}
		return nil
	})
}

func (d *Dao) do(fn func(db mongo.DB) error) error {
	return mongo.With(d.dbName, fn)
}

func create() (interface{}, error) {
	d := &Dao{
		dbName: DBName,
		logger: log.Get(PKGNAME),
	}

	e := d.Init()
	return d, e
}
