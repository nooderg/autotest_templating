package parsing

import (
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type yamlparameters struct {
	Parameters []yamlrouteParameters `yaml:"parameters"`
}

type yamlrouteParameters struct {
	Name     string     `yaml:"name"`
	In       string     `yaml:"in"`
	Required bool       `yaml:"required"`
	Value    string     `yaml:"-"`
	Schema   yamlSchema `yaml:"schema"`
}

type yamlSchema struct {
	Type   string `yaml:"type"`
	Format string `yaml:"format"`
}

func parseParameters(routeData string) ([]Parameter, error) {
	yamlparameters := yamlparameters{}
	yaml.Unmarshal([]byte(routeData), &yamlparameters)

	parameters := make([]Parameter, 0)
	for _, yamlparam := range yamlparameters.Parameters {
		parameter := Parameter{
			Name:     yamlparam.Name,
			Type:     yamlparam.Schema.Type,
			Required: yamlparam.Required,
		}
		parameters = append(parameters, parameter)
	}

	reName := regexp.MustCompile(`- name:`)
	rawYamlParameters := reName.Split(routeData, -1)[1:]
	for i, rawYamlParameter := range rawYamlParameters {
		reValue := regexp.MustCompile(`# autotest_value:.+`)
		value := reValue.FindString(rawYamlParameter)
		if len(value) != 0 {
			parameters[i].Value = strings.Split(value, "# autotest_value: ")[1]
		}
	}

	return parameters, nil
}

type yamlServers struct {
	Servers []yamlUrl `yaml:"servers"`
}

type yamlUrl struct {
	Url string `yaml:"url"`
}

func getFullUrl(yamlString string) string {
	var fullUrl yamlServers
	yaml.Unmarshal([]byte(yamlString), &fullUrl)
	return fullUrl.Servers[0].Url
}
