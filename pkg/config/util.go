package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	cenum "mkit/pkg/enum"

	"github.com/go-playground/validator/v10"
)

var (
	CoreBaseConfigPath      = "./config/core"
	DashboardBaseConfigPath = "./config/dashboard"
	BaseConfigType          = "yaml"
	ConfigName              = "config"
	EnvPrefix               = "App"
)

func LoadConfig(path, configName, configType, envPrefix string, cfg any) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.SetEnvPrefix(envPrefix)

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(cfg)
}

var configName = os.Getenv("ENV")

// getConfigs TODO: replace *App return type with suitable config type
func getConfigs(baseConfigPath string) (*App, error) {
	var (
		cfg            App
		configFileName = ConfigName
	)

	if configName == cenum.EnvironmentDevelopment.String() || configName == cenum.EnvironmentProduction.String() {
		configFileName = fmt.Sprintf("%s-%s", configFileName, configName)
	}

	if err := LoadConfig(baseConfigPath, configFileName, BaseConfigType, EnvPrefix, &cfg); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("config data error: %w", err)
	}

	return &cfg, nil
}
