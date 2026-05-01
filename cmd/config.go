package cmd

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

func loadConfig(v *viper.Viper) error {
	v.SetConfigFile(".shikai.yml")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if errors.As(err, &notFound) || os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}
