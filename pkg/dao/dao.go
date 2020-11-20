package dao

type Dao struct {
	Mysql *MysqlPool
}

func New(conf *MysqlPoolConfig) *Dao {
	return &Dao{
		Mysql: NewMysqlPool(conf),
	}
}
