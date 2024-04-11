package config

import (
	"organization/model/dto"
	"github.com/spf13/viper"
)

var ConfigVal dto.Config

func LoadConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&ConfigVal)
	return err
}
