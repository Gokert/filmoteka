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
	v.SetConfigFile(cfgPath)
	v.SetConfigType(strings.TrimPrefix(filepath.Ext(cfgPath), "."))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &DbPsxConfig{
		User:         v.GetString("user"),
		Password:     v.GetString("password"),
		Dbname:       v.GetString("dbname"),
		Host:         v.GetString("host"),
		Port:         v.GetInt("port"),
		Sslmode:      v.GetString("sslmode"),
		MaxOpenConns: v.GetInt("max_open_conns"),
		Timer:        v.GetInt("timer"),
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
