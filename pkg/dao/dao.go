package dao

import "github.com/Shanghai-Lunara/pkg/zaplogger"

type Dao struct {
	Mysql *MysqlPool
}

var d *Dao

func New(conf *MysqlPoolConfig) *Dao {
	d = &Dao{
		Mysql: NewMysqlPool(conf),
	}
	return d
}

// Get return the pointer of the Dao
func Get() *Dao {
	if d.Mysql == nil {
		zaplogger.Sugar().Fatal("error: nil Mysql, please call New() before Get()")
	}
	return d
}
