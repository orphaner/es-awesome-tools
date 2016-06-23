package at_mapping

import (
	"encoding/json"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/orphaner/es-awesome-tools/eslib"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type (
	TemplateLink struct {
		IndexName        string
		TypeName         string
		TemplateName     string
		EffectiveMapping eslib.MappingsJson
		ExpectedMapping  eslib.MappingsJson
	}
	applicableTemplate struct {
		name     string
		template eslib.TemplateJson
	}
)

var (
	esClient  *eslib.EsClient
	templates eslib.TemplateResponse
	mappings  eslib.MappingResponse
)

func FillInData(esClientParam *eslib.EsClient, filterIndex string, filterTypes string) {
	esClient = esClientParam

	templates = getTemplate()
	mappings = getMapping(filterIndex, filterTypes)
}

func getTemplate() eslib.TemplateResponse {

	request, err := esClient.NewRequest("GET", "_template", "")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	var template eslib.TemplateResponse
	json.NewDecoder(resp.Body).Decode(&template)

	return template
}

func getMapping(filterIndex string, filterTypes string) eslib.MappingResponse {

	url := fmt.Sprintf("%s/_mapping/%s", filterIndex, filterTypes)
	request, err := esClient.NewRequest("GET", url, "")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	var mapping eslib.MappingResponse
	json.NewDecoder(resp.Body).Decode(&mapping)

	return mapping
}

func GetIndexTypeAndTemplateLink() []TemplateLink {
	var links []TemplateLink

	// Parcours des index
	for indexName, indexValue := range mappings {

		// Puis des types
		for typeName, mapping := range indexValue.MappingsByType {

			expectedMapping := calculateExpectedMapping(indexName, typeName)
			links = append(links, TemplateLink{
				IndexName:        indexName,
				TypeName:         typeName,
				TemplateName:     "templateName",
				EffectiveMapping: mapping,
				ExpectedMapping:  expectedMapping,
			})
		}
	}
	sort.Sort(byIndexAndTypeSort(links))
	return links
}

func calculateExpectedMapping(indexName string, typeName string) (expectedMapping eslib.MappingsJson) {
	var applicableTemplates []applicableTemplate = searchForApplicableTemplate(indexName)
	sort.Sort(byOrderSort(applicableTemplates))

	for _, templateValue := range applicableTemplates {
		mergo.Map(&expectedMapping, templateValue.template.MappingsByType[typeName])
	}
	return expectedMapping
}

func searchForApplicableTemplate(indexName string) (result []applicableTemplate) {
	// Pour chacun des templates
	for templateName, templateValue := range templates {
		pattern := getRegexPatternFromTemplateValue(templateValue.TemplateIndexPattern)
		regex := regexp.MustCompile(pattern)

		// On regarde si le pattern correspond à l'index
		if regex.MatchString(indexName) {
			result = append(result, applicableTemplate{
				name:     templateName,
				template: templateValue,
			})
		}
	}
	return result
}

func getRegexPatternFromTemplateValue(templateIndexPattern string) (pattern string) {
	pattern = strings.Replace(templateIndexPattern, ".", "\\.", -1)
	pattern = strings.Replace(pattern, "*", ".*", -1)
	return pattern
}
