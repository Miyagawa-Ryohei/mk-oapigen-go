package usecase

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	"mk-oapigen-go/entity"
)

type ServerBuilder struct{}

func (b ServerBuilder) BuildServerSchema(spec *openapi3.Server) entity.Server {
	p, found := spec.Extensions["x-default-port"]
	pBuf, ok := p.(json.RawMessage)
	port := "5555"
	if found && ok {
		if err := json.Unmarshal(pBuf, &port); err != nil {
			port = "5555"
		}
	}

	h, ok := spec.Extensions["x-default-host"]
	hBuf, ok := h.(json.RawMessage)
	host := "0.0.0.0"
	if found && ok {
		if err := json.Unmarshal(hBuf, &host); err != nil {
			port = "0.0.0.0"
		}
	}

	prm, found := spec.Extensions["x-prometheus-exporter"]
	prmBuf, ok := prm.(json.RawMessage)
	prometheus := true
	if found && ok {
		if err := json.Unmarshal(prmBuf, &prometheus); err != nil {
			prometheus = true
		}
	}

	return entity.Server{
		Port:       port,
		Host:       host,
		Prometheus: prometheus,
	}
}
