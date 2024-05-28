package configs

import (
	"github.com/spf13/viper"
)

type DbPsxConfig struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Dbname       string `yaml:"dbname"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Sslmode      string `yaml:"sslmode"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	Timer        int    `yaml:"timer"`
}

type DbRedisCfg struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	DbNumber int    `yaml:"db"`
	Timer    int    `yaml:"timer"`
}

func GetPsxConfig() (*DbPsxConfig, error) {
	v := viper.GetViper()
	v.AutomaticEnv()

	cfg := &DbPsxConfig{
		User:         v.GetString("POSTGRES_USER"),
		Password:     v.GetString("POSTGRES_PASSWORD"),
		Dbname:       v.GetString("POSTGRES_DB"),
		Host:         v.GetString("POSTGRES_HOST"),
		Port:         v.GetInt("POSTGRES_PORT"),
		Sslmode:      v.GetString("POSTGRES_SSLMODE"),
		MaxOpenConns: v.GetInt("POSTGRES_MAXCONNS"),
		Timer:        v.GetInt("POSTGRES_TIMER"),
	}

	return cfg, nil
}

func GetRedisConfig() (*DbRedisCfg, error) {
	v := viper.GetViper()
	v.AutomaticEnv()

	cfg := &DbRedisCfg{
		Host:     v.GetString("REDIS_ADDR"),
		Password: v.GetString("REDIS_PASSWORD"),
		DbNumber: v.GetInt("REDIS_DB"),
		Timer:    v.GetInt("REDIS_TIMER"),
	}

	return cfg, nil
}
