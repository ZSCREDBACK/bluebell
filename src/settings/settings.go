package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局变量,用于保存所有的配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	StartTime    string `mapstructure:"startTime"`
	MachineID    int64  `mapstructure:"machineId"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	*LogConfig   `mapstructure:"log"`
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge"`
}

type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Dbname       string `mapstructure:"dbname"`
	MaxOpenConns int    `mapstructure:"maxOpenConns"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolSize"`
}

func Init() (err error) {
	// 加载配置文件的路径及文件名
	// viper.SetConfigName("config") // 当目录中存在多个config.xxx文件(不同扩展名),会发生空指针错误
	// viper.SetConfigType("yaml")   // 专门用于从远程服务器获取配置文件时,指定文件的格式(加载本地文件,则此项不生效)
	// 目录中存在多个配置文件时,可以用以下方法指定要加载的配置文件
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")

	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("读取配置文件失败: %v\n", err)
		return
	}

	// 将配置文件的内容解析到结构体中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("Viper unmarshal failed, err: %v\n", err)
		// 这边为什么不用return
	}

	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		// 重新反射
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("Viper unmarshal failed, err: %v\n", err)
		}
	})

	return
}
