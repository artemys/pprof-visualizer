package services

import (
	"github.com/artemys/pprof-visualizer/configs"
)

type Services struct {
}

func InitServices(config configs.ApiConfig) *Services {
	return &Services{}
}
