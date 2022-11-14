package usecase

import "testing"

func TestRouterGenerator_PrintRouter(t *testing.T) {
	generator := NewRouterGenerator("infra", "entity")
	generator.PrintRouter(Route)
}
