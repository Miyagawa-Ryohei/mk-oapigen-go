package usecase

import "github.com/getkin/kin-openapi/openapi3"

func ReadSpec(p string) (*openapi3.T, error) {
	sp, err := openapi3.NewLoader().LoadFromFile(p)
	if err != nil {
		return nil, err
	}
	return sp, nil
}
