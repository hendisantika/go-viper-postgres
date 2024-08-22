package main

import (
	"github.com/spf13/viper"
	"reflect"
)

func LoadConfig[TConfig any](path string) (cfg TConfig, err error) {
	// Read app.env from specified directory
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// try read cfg file
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	// try load settings from env vars

	// Retrieve the underlying type of variable `i`.
	r := reflect.TypeOf(cfg)

	// Iterate over each field for the type
	for j := 0; j < r.NumField(); j++ {
		f := r.Field(j)

		// Bind the environment variable.
		if err = viper.BindEnv(f.Name); err != nil {
			return
		}
	}

	if err = viper.Unmarshal(&cfg); err != nil {
		return
	}

	return
}
