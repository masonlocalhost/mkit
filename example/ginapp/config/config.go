package config

import (
	"fmt"
	"mkit/pkg/config"
	cenum "mkit/pkg/enum"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var (
	BaseConfigPath = "./config"
	BaseConfigType = "yaml"
	ConfigName     = "config"
	EnvPrefix      = "WEBAPP"
)

type Config struct {
	config.App `mapstructure:",squash"`
}

var configName = os.Getenv("ENV")

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

func getConfigs(baseConfigPath string) (*Config, error) {
	var (
		cfg            Config
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

func GetConfig() (*Config, error) {
	cfg, err := getConfigs(BaseConfigPath)
	if err != nil {
		return nil, err
	}

	if cfg.HTTP == nil {
		return nil, fmt.Errorf("'http' config is required")
	}

	return cfg, nil
}
