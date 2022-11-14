package usecase

import (
	"github.com/getkin/kin-openapi/openapi3"
	"mk-oapigen-go/entity"
)

type ServerBuilder struct{}

func (b ServerBuilder) BuildServerSchema(spec *openapi3.Server) entity.Server {
	p, found := spec.Extensions["x-default-port"]
	port, ok := p.(string)
	if (!found) || (!ok) {
		port = "5555"
	}

	h, ok := spec.Extensions["x-default-host"]
	host, ok := h.(string)
	if !ok {
		host = "0.0.0.0"
	}

	return entity.Server{
		Port: port,
		Host: host,
	}
}
