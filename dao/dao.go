package dao

import (
	"errors"
	"time"
	"github.com/cuigh/auxo/config"
	"github.com/cuigh/auxo/util/lazy"
	"github.com/cuigh/auxo/log"
	"github.com/globalsign/mgo"
	"github.com/dongdongking008/softtrans/misc"
)

var (
	indexes = map[string][]mgo.Index{
		"transaction": {
			mgo.Index{Key: []string{"trans_id"}, Unique: true},
			mgo.Index{Key: []string{"status"}, },
			mgo.Index{Key: []string{"enter_time"}, },
			mgo.Index{Key: []string{"expire_time"}, },
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

type database struct {
	db *mgo.Database
}

func (d *database) Close() {
	d.db.Session.Close()
}

func (d *database) C(name string) *mgo.Collection {
	return d.db.C(name)
}

func (d *database) Run(cmd, result interface{}) error {
	return d.db.Run(cmd, result)
}

type Dao struct {
	logger  log.Logger
	session *mgo.Session
}

func New(addr string) (*Dao, error) {
	if addr == "" {
		return nil, errors.New("database address must be configured for mongo storage")
	}

	s, err := mgo.DialWithTimeout(addr, time.Second*5)
	if err != nil {
		return nil, err
	}

	d := &Dao{
		session: s,
		logger:  log.Get("mongo"),
	}
	return d, nil
}

func (d *Dao) Init() {
	db := d.db()
	defer db.Close()

	for name, ins := range indexes {
		c := db.C(name)
		for _, in := range ins {
			err := c.EnsureIndex(in)
			if err != nil {
				d.logger.Warnf("Ensure index %s-%v failed: %v", name, in.Key, err)
			}
		}
	}
}

func (d *Dao) Close() {
	d.session.Close()
}

func (d *Dao) db() *database {
	return &database{
		db: d.session.Copy().DB(""),
	}
}

func (d *Dao) do(fn func(db *database)) {
	db := d.db()
	defer db.Close()

	fn(db)
}

func create() (interface{}, error) {
	d, err := New(config.GetString(misc.KeyDBAddress))

	if err == nil {
		d.Init()
	}
	return d, err
}