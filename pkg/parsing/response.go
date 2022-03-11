package parsing

import (
	"gopkg.in/yaml.v3"
)

type yamlResponses struct {
	Response map[interface{}]*yamlResponse `yaml:"responses"`
}

type yamlResponse struct {
	Description string `yaml:"description"`
}

func parseResponses(routeData string) (Response, error) {
	yamlResponses := yamlResponses{}
	yaml.Unmarshal([]byte(routeData), &yamlResponses)

	response := Response{}
	if yamlResponses.Response[200] != nil {
		response.Status = 200
	} else if yamlResponses.Response[201] != nil {
		response.Status = 201
	} else if yamlResponses.Response[202] != nil {
		response.Status = 202
	} else if yamlResponses.Response[203] != nil {
		response.Status = 203
	} else if yamlResponses.Response[204] != nil {
		response.Status = 204
	}

	return response, nil
}
