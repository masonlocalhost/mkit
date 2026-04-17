package config

import (
	"mkit/pkg/enum"
	"time"
)

type App struct {
	Environment enum.Environment `mapstructure:"environment" validate:"required"`
	Timezone    string           `mapstructure:"timezone" validate:"required"`
	Version     string           `mapstructure:"version" validate:"required"`
	Postgres    *Postgres        `mapstructure:"postgres"`
	Redis       *Redis           `mapstructure:"redis"`
	HTTP        *HTTP            `mapstructure:"http"`
	GRPC        *GRPC            `mapstructure:"grpc"`
	Connect     *Connect         `mapstructure:"connect"`
	Log         *Log             `mapstructure:"log"`
	Tracing     Tracing          `mapstructure:"tracing"`
	Swagger     *Swagger         `mapstructure:"swagger"`
	Minio       *Minio           `mapstructure:"minio"`
	RabbitMQ    *RabbitMQ        `mapstructure:"rabbitmq"`
}

type Postgres struct {
	ConnectionParams *PostgresConnectionParams `mapstructure:"connectionParams" validate:"required"`
	IsMigrateSchema  bool                      `mapstructure:"isMigrateSchema"`
	MaxOpenConn      int                       `mapstructure:"maxOpenConn" validate:"required,gt=0"`
	MaxIdleConn      int                       `mapstructure:"maxIdleConn" validate:"required,gt=0"`
}

type PostgresConnectionParams struct {
	Host     string `mapstructure:"host" validate:"required,hostname|ip"`
	Port     int    `mapstructure:"port" validate:"required,gt=0"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	DBName   string `mapstructure:"dbName" validate:"required"`
	SSLMode  string `mapstructure:"sslMode" validate:"required,oneof=disable require verify-full"`
}

type Redis struct {
	Host     string `mapstructure:"host" validate:"required,hostname|ip"`
	Port     string `mapstructure:"port,gt=0"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type GRPC struct {
	Host                 string `mapstructure:"host" validate:"required,hostname|ip"`
	Port                 string `mapstructure:"port" validate:"required,gt=0"`
	MaxConnectionAge     string `mapstructure:"maxConnectionAge"`
	ReflectionEnabled    bool   `mapstructure:"reflectionEnabled"`
	JsonTranscodeEnabled bool   `mapstructure:"jsonTranscodeEnabled"`
}

type HTTP struct {
	Host string `mapstructure:"host" validate:"required,hostname|ip"`
	Port string `mapstructure:"port" validate:"required,gt=0"`
}

type Connect struct {
	Host string `mapstructure:"host" validate:"required,hostname|ip"`
	Port string `mapstructure:"port" validate:"required,gt=0"`
}

type Tracing struct {
	ServiceName   string `mapstructure:"serviceName"`
	Enabled       bool   `mapstructure:"enabled"`
	CollectorHost string `mapstructure:"collectorHost"`
	CollectorPort int    `mapstructure:"collectorPort"`
}

type Log struct {
	Level       string `mapstructure:"level" validate:"required,oneof=panic fatal error warn info debug trace"`
	LogFilePath string `mapstructure:"logFilePath" validate:"required"`
}

type Swagger struct {
	Path string `mapstructure:"path" validate:"required"`
}

type GRPCClient struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Minio struct {
	Host       string `mapstructure:"host" validate:"required,hostname|ip"`
	Port       int    `mapstructure:"port" validate:"required,gt=0"`
	AccessKey  string `mapstructure:"accessKey" validate:"required"`
	SecretKey  string `mapstructure:"secretKey" validate:"required"`
	BucketName string `mapstructure:"bucketName" validate:"required"`
	SSLEnabled bool   `mapstructure:"sslEnabled"`
}

type RabbitMQ struct {
	Host     string `mapstructure:"host" validate:"required,hostname|ip"`
	Port     int    `mapstructure:"port" validate:"required,gt=0"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
}

type Cronjob struct {
	ID          string        `mapstructure:"id"`
	Spec        string        `mapstructure:"spec"`
	TaskTimeout time.Duration `mapstructure:"taskTimeout"`
	Disabled    bool          `mapstructure:"disabled"`
}
