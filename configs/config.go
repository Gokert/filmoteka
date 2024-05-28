package configs

import (
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
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

func GetPsxConfig(cfgPath string) (*DbPsxConfig, error) {
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

func GetRedisConfig(cfgPath string) (*DbRedisCfg, error) {
	v := viper.GetViper()
	v.SetConfigFile(cfgPath)
	v.SetConfigType(strings.TrimPrefix(filepath.Ext(cfgPath), "."))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &DbRedisCfg{
		Host:     v.GetString("host"),
		Password: v.GetString("password"),
		DbNumber: v.GetInt("db"),
		Timer:    v.GetInt("timer"),
	}

	return cfg, nil
}
