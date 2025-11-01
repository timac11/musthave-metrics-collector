package handler

import (
	"github.com/timac11/musthave-metrics-collector/internal/handler/templates"
	"github.com/timac11/musthave-metrics-collector/internal/service"
	"html/template"
)

type ApplicationAPIContainer struct {
	service      service.Service
	templatesMap map[string]*template.Template
}

func NewApplicationAPIContainer(s service.Service) *ApplicationAPIContainer {
	indexTpl, err := template.ParseFS(templates.Index, "index.gohtml")
	templatesMap := make(map[string]*template.Template)

	if err == nil {
		templatesMap["index"] = indexTpl
	}

	container := &ApplicationAPIContainer{
		service:      s,
		templatesMap: templatesMap,
	}

	return container
}
