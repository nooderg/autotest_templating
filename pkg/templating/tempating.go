package templating

import (
	"bytes"
	"html/template"

	"github.com/nooderg/autotest_templating/pkg/parsing"
)

const TAVERN_TEMPLATE_BASE = `
---
test_name: Autotest templating
stages: 
  {{range .}}
  - name: {{ .Path}}
    delay_before: 5
    max_retries: 5
    request:
      url: {{ .Path}}
      method: {{ .Method}}
    response:
      status_code: {{ .Status}}
  {{end}}
`

type TemplatingData struct {
	Path   string
	Method string
	Status int
}

func TemplateFile(pathDatas []parsing.PathData) (*bytes.Buffer, error) {
	tdatas := make([]TemplatingData, 0)
	for _, pathData := range pathDatas {
		for _, routeData := range pathData.RouteData {
			tdata := TemplatingData{
				Path:   pathData.BaseUrl + routeData.Path,
				Method: routeData.Method,
				Status: int(routeData.Response.Status),
			}

			tdatas = append(tdatas, tdata)
		}
	}

	t, err := template.New("TAVERN_TEMPLATE").Parse(TAVERN_TEMPLATE_BASE)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	err = t.Execute(&b, tdatas)
	if err != nil {
		return nil, err
	}
	return &b, nil
}
