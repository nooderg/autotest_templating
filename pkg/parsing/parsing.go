package parsing

import (
	"io/ioutil"
	"log"
	"regexp"
)

type Parameter struct {
	Name     string
	Type     string
	Required bool
}

type Response struct {
	Status uint16
}


type RouteData struct {
	Method     string
	Summary    string
	Parameters []Parameter
	Response   Response
}

type PathData struct {
	Path      string
	RouteData []RouteData
}



type preprocessedData struct {
	Path string
	Data string
}



func ParseOpenApi() (*PathData, error) {
	yamlFile, err := ioutil.ReadFile("./example/oa-example.yaml")
	if err != nil {
		return nil, err
	}

	_, rawishDatas := isolateComponentsAndRoutes(string(yamlFile))

	log.Println(rawishDatas[0])


	var data PathData

	return &data, nil
}

func isolateComponentsAndRoutes(yamlString string) (string, []preprocessedData) {
	reComp := regexp.MustCompile(`components:`)

	components := reComp.Split(string(yamlString), -1)[1]

	reRoutes := regexp.MustCompile(`[ \s]{2}\/.+`)
	urls := reRoutes.FindAllString(reComp.Split(string(yamlString), -1)[0], -1)
	routes := reRoutes.Split(reComp.Split(string(yamlString), -1)[0], -1)[1:]


	rawishDatas := make([]preprocessedData, 0)
	for i, _ := range urls {
		rawishData := preprocessedData{
			Path: urls[i], 
			Data: routes[i],
		}
		rawishDatas = append(rawishDatas, rawishData)
	}

	return components, rawishDatas
}

func splitRoutesMethod(route string) []string {
	reMethods := regexp.MustCompile(`(post:|get:)`)
	return reMethods.Split(route, -1)
}

func getRouteMethodData(route string) {

}
