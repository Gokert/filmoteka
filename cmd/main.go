package main

import (
	"filmoteka/configs"
	"filmoteka/configs/logger"
	delivery "filmoteka/delivery/http"
	"filmoteka/usecase"
	"github.com/joho/godotenv"
	_ "github.com/swaggo/swag"
)

// @title filmoteka App API
// @version 1.0
// @description API Server fot Application
// @host 127.0.0.1:8081
// @BasePath /
func main() {
	log := logger.GetLogger()
	err := godotenv.Load()
	if err != nil {
		log.Error("load .env error: ", err)
		return
	}

	psxCfg, err := configs.GetPsxConfig()
	if err != nil {
		log.Error("Create psx config error: ", err)
		return
	}

	redisCfg, err := configs.GetRedisConfig()
	if err != nil {
		log.Error("Create redis config error: ", err)
		return
	}

	core, err := usecase.GetCore(psxCfg, redisCfg, log)
	if err != nil {
		log.Error("Create core error: ", err)
		return
	}

	api := delivery.GetApi(core, log)

	log.Info("Server running")
	err = api.ListenAndServe("8081")
	if err != nil {
		log.Error("ListenAndServe error: ", err)
		return
	}

}
