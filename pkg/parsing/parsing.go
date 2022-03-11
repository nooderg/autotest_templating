package parsing

import (
	"regexp"
)

type Parameter struct {
	Name     string
	Type     string
	Value    string
	Required bool
}

type Response struct {
	Status uint16
}

type RouteData struct {
	Path       string
	Method     string
	Summary    string
	Parameters []Parameter
	Response   Response
}

type PathData struct {
	BaseUrl   string
	RouteData []RouteData
}

type preprocessedData struct {
	Path    string
	RawData string
}

func ParseOpenApi(yamlFile string) ([]PathData, error) {
	_, rawishDatas := preprocessData(string(yamlFile))

	PathDatas := make([]PathData, len(rawishDatas))
	for i, rawishData := range rawishDatas {
		PathDatas[i] = processRawData(rawishData)
		PathDatas[i].BaseUrl = getFullUrl(string(yamlFile))

		for k, route := range PathDatas[i].RouteData {
			rePath := regexp.MustCompile("{[a-zA-Z]+}")
			reMatch := rePath.FindAllString(route.Path, -1)
			if len(reMatch) != 0 {
				for _, param := range route.Parameters {
					if param.Name == reMatch[0][1:len(reMatch[0])-1] {
						newUrl := rePath.ReplaceAll([]byte(route.Path), []byte(param.Value))
						PathDatas[i].RouteData[k].Path = string(newUrl)
					}
				}
			}
		}
	}

	return PathDatas, nil
}

func preprocessData(yamlString string) (string, []preprocessedData) {
	reComp := regexp.MustCompile(`components:`)

	paths := reComp.Split(string(yamlString), -1)[0]
	components := reComp.Split(string(yamlString), -1)[1]

	reRoutes := regexp.MustCompile(`[ \s]{2}\/.+`)
	urls := reRoutes.FindAllString(paths, -1)
	routes := reRoutes.Split(paths, -1)[1:]

	rawishDatas := make([]preprocessedData, 0)
	for i, _ := range urls {
		rawishData := preprocessedData{
			Path:    urls[i],
			RawData: routes[i],
		}
		rawishDatas = append(rawishDatas, rawishData)
	}

	return components, rawishDatas
}

func processRawData(data preprocessedData) PathData {
	reMethods := regexp.MustCompile(`(post:|get:|put:|patch:|options:)`)
	methods := reMethods.FindAllString(data.RawData, -1)
	routes := reMethods.Split(data.RawData, -1)[1:]

	routesData := make([]RouteData, 0)
	for i, _ := range methods {
		parameters, _ := parseParameters(routes[i])
		response, _ := parseResponses(routes[i])
		routeData := RouteData{
			Method:     methods[i][:len(methods[i])-1],
			Parameters: parameters,
			Response:   response,
			Path:       data.Path[2 : len(data.Path)-1],
		}

		routesData = append(routesData, routeData)
	}

	PathData := PathData{
		RouteData: routesData,
	}

	return PathData
}
