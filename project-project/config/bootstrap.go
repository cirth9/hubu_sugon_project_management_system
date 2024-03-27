package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
)

//var BC = InitBootstrap()

type BootConf struct {
	viper       *viper.Viper
	NacosConfig *NacosConfig
}

func (c *BootConf) ReadNacosConfig() {
	nc := &NacosConfig{}
	_ = c.viper.UnmarshalKey("nacos", nc)
	zap.S().Info("nacos config: ", nc)
	c.NacosConfig = nc
}

type NacosConfig struct {
	Namespace   string
	Group       string
	IpAddr      string
	Port        int
	ContextPath string
	Scheme      string
}

func InitBootstrap() *BootConf {
	conf := &BootConf{viper: viper.New()}
	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("bootstrap")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath(workDir + "/config")
	conf.viper.AddConfigPath("D:\\项目实战\\project_management_system\\project-project\\config")
	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	conf.ReadNacosConfig()
	return conf
}
