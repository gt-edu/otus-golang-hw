package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

func NewConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error during reading configuration: %v", err)
		return nil, err
	}
	conf := &Config{}
	err := viper.Unmarshal(conf)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error during unmarshalling configuration: %v", err)
		return nil, err
	}

	return conf, nil
}

// TODO
