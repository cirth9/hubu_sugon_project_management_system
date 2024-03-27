package gorms

import (
	"context"
	"gorm.io/gorm"
)

var _db *gorm.DB

//func init() {
//	if config.C.DbConfig.Separation {
//		//开启读写分离
//		username := config.C.DbConfig.Master.Username //账号
//		password := config.C.DbConfig.Master.Password //密码
//		host := config.C.DbConfig.Master.Host         //数据库地址，可以是Ip或者域名
//		port := config.C.DbConfig.Master.Port         //数据库端口
//		Dbname := config.C.DbConfig.Master.Db         //数据库名
//		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
//		var err error
//		_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//			Logger: logger.Default.LogMode(logger.Info),
//		})
//		if err != nil {
//			panic("连接数据库失败, error=" + err.Error())
//		}
//		//slave
//		replicas := []gorm.Dialector{}
//		for _, v := range config.C.DbConfig.Slave {
//			username := v.Username //账号
//			password := v.Password //密码
//			host := v.Host         //数据库地址，可以是Ip或者域名
//			port := v.Port         //数据库端口
//			Dbname := v.Db         //数据库名
//			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
//			cfg := mysql.Config{
//				DSN: dsn,
//			}
//			replicas = append(replicas, mysql.New(cfg))
//		}
//		_db.Use(dbresolver.Register(dbresolver.Config{
//			Sources: []gorm.Dialector{mysql.New(mysql.Config{
//				DSN: dsn,
//			})},
//			Replicas: replicas,
//			Policy:   dbresolver.RandomPolicy{},
//		}).SetMaxOpenConns(200).SetMaxIdleConns(10))
//	} else {
//		//配置MySQL连接参数
//		username := config.C.MysqlConfig.Username //账号
//		password := config.C.MysqlConfig.Password //密码
//		host := config.C.MysqlConfig.Host         //数据库地址，可以是Ip或者域名
//		port := config.C.MysqlConfig.Port         //数据库端口
//		Dbname := config.C.MysqlConfig.Db         //数据库名
//		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, Dbname)
//		var err error
//		_db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//			Logger: logger.Default.LogMode(logger.Info),
//		})
//		if err != nil {
//			panic("连接数据库失败, error=" + err.Error())
//		}
//	}
//}

func GetDB() *gorm.DB {
	return _db
}

func SetDB(db *gorm.DB) {
	_db = db
}

type GormConn struct {
	db *gorm.DB
	tx *gorm.DB
}

func (g *GormConn) Begin() {
	g.tx = GetDB().Begin()
}

func New() *GormConn {
	return &GormConn{db: GetDB()}
}

func NewTran() *GormConn {
	return &GormConn{db: GetDB(), tx: GetDB()}
}

func (g *GormConn) Session(ctx context.Context) *gorm.DB {
	return g.db.Session(&gorm.Session{Context: ctx})
}

func (g *GormConn) Rollback() {
	g.tx.Rollback()
}
func (g *GormConn) Commit() {
	g.tx.Commit()
}

func (g *GormConn) Tx(ctx context.Context) *gorm.DB {
	return g.tx.WithContext(ctx)
}
