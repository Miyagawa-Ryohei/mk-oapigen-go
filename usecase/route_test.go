package usecase

import (
	"mk-oapigen-go/entity"
	"testing"
)

var Route = []entity.Route{
	{
		Name: "AAA",
		Path: "/aaa",
		Methods: []entity.Method{
			{
				Type:     "GET",
				Name:     "GetAAA",
				Request:  TestRequest1,
				Response: TestRequest2,
			}, {
				Type:     "POST",
				Name:     "PostAAA",
				Request:  TestRequest1,
				Response: TestRequest2,
			},
		},
	}, {
		Name: "AAB",
		Path: "/aab",
		Methods: []entity.Method{
			{
				Type:     "GET",
				Name:     "GetAAB",
				Request:  TestRequest1,
				Response: TestRequest2,
			}, {
				Type:     "DELETE",
				Name:     "DeleteAAB",
				Request:  TestRequest1,
				Response: TestRequest2,
			},
		},
	}, {
		Name: "AAC",
		Path: "/aac",
		Methods: []entity.Method{
			{
				Type:     "GET",
				Name:     "GetAAC",
				Request:  TestRequest1,
				Response: TestRequest2,
			}, {
				Type:     "PUT",
				Name:     "PutAAC",
				Request:  TestRequest3,
				Response: TestRequest1,
			},
		},
	},
}

func TestRouteGenerator_PrintRoute(t *testing.T) {
	generator := NewRouteGenerator("entity", "")
	generator.PrintRoute(Route)
}
