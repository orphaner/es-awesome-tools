package at_mapping

import (
	"fmt"
	"github.com/orphaner/es-awesome-tools/eslib"
	"net/http"
	"testing"
)

const (
	startupPort = "7889"
)

type httpMock struct {
	mappingJson  string
	templateJson string
}

func (httpMock *httpMock) setMapping(mapping string) {
	httpMock.mappingJson = mapping
}

func (httpMock *httpMock) setTemplate(template string) {
	httpMock.templateJson = template
}

func (httpMock *httpMock) startMockServer() {
	http.HandleFunc("/_template",
		func(responseWriter http.ResponseWriter, request *http.Request) {
			fmt.Fprint(responseWriter, httpMock.templateJson)
		})
	http.HandleFunc("/*/_mapping/*",
		func(responseWriter http.ResponseWriter, request *http.Request) {
			fmt.Fprint(responseWriter, httpMock.mappingJson)
		})
	http.ListenAndServe(":" + startupPort, nil)
}

type mappingTest struct {
	httpMock    httpMock
	mappingTool *MappingTool
}

func (test *mappingTest) TestMain(m *testing.M) {
	go test.httpMock.startMockServer()

	esClient := eslib.NewEsClient()
	esClient.SetFromFlag("http://localhost:" + startupPort)
	test.mappingTool = NewMappingTool(&esClient)

	m.Run()
}

func (test *mappingTest) TestMockMechanism(t *testing.T) {
	test.httpMock.setMapping(`{"mappingKey": {}}`)
	test.httpMock.setTemplate(`{"templateKey": {}}`)

	mappingTool.FillInData("*", "*")
	if _, ok := mappingTool.mappings["mappingKey"]; !ok {
		t.Errorf("No mapping found for index 'mappingKey'")
	}
	if _, ok := mappingTool.Templates["templateKey"]; !ok {
		t.Errorf("No template 'templateKey' found")
	}
}
