package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"k8s.io/klog/v2"
	"time"
)

const (
	DBDSNFormat = "%s:%s@tcp(%s:%d)/%s?timeout=5s&readTimeout=5s&writeTimeout=5s&parseTime=true&loc=Local&charset=utf8mb4,utf8"
)

type MysqlConfig struct {
	Host            string `json:"host" yaml:"host"`
	Port            int    `json:"port" yaml:"port"`
	User            string `json:"user" yaml:"user"`
	Password        string `json:"password" yaml:"password"`
	Database        string `json:"database" yaml:"database"`
	MaxIdleConns    int    `json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int    `json:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
}

type MysqlPoolConfig struct {
	Master MysqlConfig `yaml:"master,flow"`
	Slave  MysqlConfig `yaml:"slave,flow"`
}

type MysqlPool struct {
	conf  *MysqlPoolConfig
	master *sql.DB
	slave *sql.DB
}

func NewMysqlPool(conf *MysqlPoolConfig) *MysqlPool {
	m := &MysqlPool{
		conf:  conf,
		master: &sql.DB{},
		slave: &sql.DB{},
	}
	// Master
	dsn := fmt.Sprintf(DBDSNFormat, conf.Master.User, conf.Master.Password, conf.Master.Host, conf.Master.Port, conf.Master.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		klog.Fatal("error Master connecting: %s", err.Error())
		return nil
	}
	db.SetMaxIdleConns(conf.Master.MaxIdleConns)
	db.SetMaxOpenConns(conf.Master.MaxOpenConns)
	db.SetConnMaxLifetime(time.Second * time.Duration(conf.Master.ConnMaxLifetime))
	m.master = db
	// Slave
	dsn2 := fmt.Sprintf(DBDSNFormat, conf.Slave.User, conf.Slave.Password, conf.Slave.Host, conf.Slave.Port, conf.Slave.Database)
	db2, err := sql.Open("mysql", dsn2)
	if err != nil {
		klog.Fatal("error Slave connecting: %s", err.Error())
		return nil
	}
	db2.SetMaxIdleConns(conf.Slave.MaxIdleConns)
	db2.SetMaxOpenConns(conf.Slave.MaxOpenConns)
	db2.SetConnMaxLifetime(time.Second * time.Duration(conf.Slave.ConnMaxLifetime))
	m.master = db
	return m
}

func (mp *MysqlPool) Master() *sql.DB {
	return mp.master
}

func (mp *MysqlPool) Slave() *sql.DB {
	return mp.slave
}
