package usecase

import (
	"mk-oapigen-go/entity"
	"testing"
)

var Server = entity.Server{
	Host: "",
	Port: "5432",
}

func TestServerGenerator_PrintRoute(t *testing.T) {
	generator := NewServerGenerator("infra", "entity")
	generator.PrintServer(Server, Route)
}
