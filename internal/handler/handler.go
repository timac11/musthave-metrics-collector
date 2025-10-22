package handler

import (
	service "github.com/timac11/musthave-metrics-collector/internal/service"
)

type ApplicationAPIContainer struct {
	service service.Service
}

func NewApplicationAPIContainer(s service.Service) *ApplicationAPIContainer {
	container := &ApplicationAPIContainer{
		service: s,
	}

	return container
}
