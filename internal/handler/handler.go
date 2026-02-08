package handler

import (
	"html/template"

	"github.com/timac11/musthave-metrics-collector/internal/audit"
	"github.com/timac11/musthave-metrics-collector/internal/handler/templates"
	"github.com/timac11/musthave-metrics-collector/internal/service"
)

type ApplicationAPIContainer struct {
	service      service.Service
	auditor      audit.Auditor
	templatesMap map[string]*template.Template
}

func NewApplicationAPIContainer(s service.Service, a audit.Auditor) *ApplicationAPIContainer {
	indexTpl, err := template.ParseFS(templates.Index, "index.gohtml")
	templatesMap := make(map[string]*template.Template)

	if err == nil {
		templatesMap["index"] = indexTpl
	}

	container := &ApplicationAPIContainer{
		service:      s,
		templatesMap: templatesMap,
		auditor:      a,
	}

	return container
}
