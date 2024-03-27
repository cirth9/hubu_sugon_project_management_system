package config

import (
	"bytes"
	"github.com/go-redis/redis/v8"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"log"
	"os"
	"test.com/project-common/logs"
)

// C todo main之前init config
var C = InitConfig()

type Config struct {
	viper        *viper.Viper
	SC           *ServerConfig
	GC           *GrpcConfig
	EtcdConfig   *EtcdConfig
	MysqlConfig  *MysqlConfig
	JwtConfig    *JwtConfig
	DbConfig     DbConfig
	JaegerConfig *JaegerConfig
}

type JaegerConfig struct {
	Endpoints string
}

type ServerConfig struct {
	Name string
	Addr string
}

type GrpcConfig struct {
	Name     string
	EtcdAddr string
	Addr     string
	Version  string
	Weight   int64
}

type EtcdConfig struct {
	Addrs []string
}

type MysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Db       string
	Name     string
}

type DbConfig struct {
	Master     MysqlConfig
	Slave      []MysqlConfig
	Separation bool
}

type JwtConfig struct {
	AccessExp     int64
	RefreshExp    int64
	AccessSecret  string
	RefreshSecret string
}

func InitConfig() *Config {
	//先读取到viper struct自己的config map内，后面直接通过get拿就好了
	conf := &Config{viper: viper.New()}
	//先从nacos读取配置，如果读取不到 在本地读取
	nacosClient := InitNacosClient()
	configYaml, err2 := nacosClient.confClient.GetConfig(vo.ConfigParam{
		DataId: "config.yaml",
		Group:  nacosClient.group,
	})
	if err2 != nil {
		log.Fatalln(err2)
	}
	err2 = nacosClient.confClient.ListenConfig(vo.ConfigParam{
		DataId: "config.yaml",
		Group:  nacosClient.group,
		OnChange: func(namespace, group, dataId, data string) {
			//
			log.Printf("load nacos config changed %s \n", data)
			err := conf.viper.ReadConfig(bytes.NewBuffer([]byte(data)))
			if err != nil {
				log.Printf("load nacos config changed err : %s \n", err.Error())
			}
			//所有的配置应该重新读取
			conf.ReLoadAllConfig()
		},
	})
	if err2 != nil {
		log.Fatalln(err2)
	}
	conf.viper.SetConfigType("yaml")
	if configYaml != "" {
		//todo 如果在nacos内读取到，就使用nacos的配置
		err := conf.viper.ReadConfig(bytes.NewBuffer([]byte(configYaml)))
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		//todo 如果没有就读取本地配置
		workDir, _ := os.Getwd()
		conf.viper.SetConfigName("config")
		conf.viper.AddConfigPath(workDir + "/config")
		err := conf.viper.ReadInConfig()
		if err != nil {
			log.Fatalln(err)
		}
	}
	conf.ReLoadAllConfig()
	return conf
}

func (c *Config) ReLoadAllConfig() {
	c.ReadServerConfig()
	c.InitZapLog()
	c.ReadGrpcConfig()
	c.ReadEtcdConfig()
	c.InitMysqlConfig()
	c.InitJwtConfig()
	c.InitDbConfig()
	c.InitJaegerConfig()
	//重新创建相关的客户端
	c.ReConnRedis()
	c.ReConnMysql()
}
func (c *Config) InitJaegerConfig() {
	mc := &JaegerConfig{
		Endpoints: c.viper.GetString("jaeger.endpoints"),
	}
	c.JaegerConfig = mc
}
func (c *Config) InitZapLog() {
	//从配置中读取日志配置，初始化日志
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxSize:       c.viper.GetInt("maxSize"),
		MaxAge:        c.viper.GetInt("maxAge"),
		MaxBackups:    c.viper.GetInt("maxBackups"),
	}
	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Config) ReadServerConfig() {
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.SC = sc
}

func (c *Config) ReadGrpcConfig() {
	gc := &GrpcConfig{}
	gc.Name = c.viper.GetString("grpc.name")
	gc.Addr = c.viper.GetString("grpc.addr")
	gc.Version = c.viper.GetString("grpc.version")
	gc.Weight = c.viper.GetInt64("grpc.weight")
	gc.EtcdAddr = c.viper.GetString("grpc.etcdAddr")
	c.GC = gc
}

func (c *Config) ReadRedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       c.viper.GetInt("redis.db"),
	}
}

func (c *Config) ReadEtcdConfig() {
	ec := &EtcdConfig{}
	var addrs []string
	err := c.viper.UnmarshalKey("etcd.addrs", &addrs)
	if err != nil {
		log.Fatalln(err)
	}
	ec.Addrs = addrs
	c.EtcdConfig = ec
}
func (c *Config) InitMysqlConfig() {
	mc := &MysqlConfig{
		Username: c.viper.GetString("mysql.username"),
		Password: c.viper.GetString("mysql.password"),
		Host:     c.viper.GetString("mysql.host"),
		Port:     c.viper.GetInt("mysql.port"),
		Db:       c.viper.GetString("mysql.db"),
	}
	c.MysqlConfig = mc
}
func (c *Config) InitJwtConfig() {
	mc := &JwtConfig{
		AccessSecret:  c.viper.GetString("jwt.accessSecret"),
		AccessExp:     c.viper.GetInt64("jwt.accessExp"),
		RefreshExp:    c.viper.GetInt64("jwt.refreshExp"),
		RefreshSecret: c.viper.GetString("jwt.refreshSecret"),
	}
	c.JwtConfig = mc
}

func (c *Config) InitDbConfig() {
	mc := DbConfig{}
	mc.Separation = c.viper.GetBool("db.separation")
	var slaves []MysqlConfig
	err := c.viper.UnmarshalKey("db.slave", &slaves)
	if err != nil {
		panic(err)
	}
	master := MysqlConfig{
		Username: c.viper.GetString("db.master.username"),
		Password: c.viper.GetString("db.master.password"),
		Host:     c.viper.GetString("db.master.host"),
		Port:     c.viper.GetInt("db.master.port"),
		Db:       c.viper.GetString("db.master.db"),
	}
	mc.Master = master
	mc.Slave = slaves
	c.DbConfig = mc
}
