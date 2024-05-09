package config

import "time"

type ApplicationConfig struct {
	Port     PortConfig     `json:"port"`
	Database DatabaseConfig `json:"database"`
	JWT      JWTConfig      `json:"jwt"`
	Encrypt  EncryptConfig  `json:"encrypt"`
}

type PortConfig struct {
	Service        int           `json:"service"`
	ServiceTimeout time.Duration `json:"servicetimeout"`
	BasePath       string        `json:"basepath"`
}

type DatabaseConfig struct {
	DBName   string        `json:"dbname"`
	Host     string        `json:"host"`
	Port     string        `json:"port"`
	User     string        `json:"user"`
	Password string        `json:"password"`
	Timeout  time.Duration `json:"timeout"`
}

type JWTConfig struct {
	AccessTokenDuration  time.Duration `json:"access_token_duration"`
	RefreshTokenDuration time.Duration `json:"refresh_token_duration"`
	SecretKey            string        `json:"secret_key"`
}

type EncryptConfig struct {
	SecretKey string `json:"secretkey"`
}
