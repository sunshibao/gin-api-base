package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Init 加载配置文件
func Init() {
	env := os.Getenv("GO_ENV")

	viper.AddConfigPath("config")
	viper.SetConfigType("json")

	// 根据环境加载不同配置
	switch env {
	case "prod":
		viper.SetConfigName("config.prod")
	default:
		viper.SetConfigName("config")
	}

	// 支持环境变量覆盖
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	log.Printf("使用配置文件: %s", viper.ConfigFileUsed())
}
