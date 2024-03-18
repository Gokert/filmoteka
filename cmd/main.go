package main

import (
	"filmoteka/configs"
	"filmoteka/configs/logger"
	delivery "filmoteka/delivery/http"
	"filmoteka/usecase"
	_ "github.com/swaggo/swag"
)

// @title filmoteka App API
// @version 1.0
// @description API Server fot Application
// @host 127.0.0.1:8081
// @BasePath /

func main() {
	log := logger.GetLogger()

	psxCfg, err := configs.GetPsxConfig("configs/db_psx.yaml")
	if err != nil {
		log.Error("Create psx config error: ", err)
		return
	}

	redisCfg, err := configs.GetRedisConfig("configs/db_redis.yaml")
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
