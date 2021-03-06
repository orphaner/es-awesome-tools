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
	MappingTool struct {
		esClient  *eslib.EsClient
		Templates eslib.TemplateResponse
		mappings  eslib.MappingResponse
	}

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

func NewMappingTool(esClientParam eslib.EsClient) *MappingTool {
	return &MappingTool{
		esClient: &esClientParam,
	}
}

func (tool *MappingTool) FillInData(filterIndex string, filterTypes string) {
	tool.Templates = tool.getTemplate()
	tool.mappings = tool.getMapping(filterIndex, filterTypes)
}

func (tool *MappingTool) getTemplate() eslib.TemplateResponse {

	request, err := (*tool.esClient).NewRequest("GET", "_template", "")
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

func (tool *MappingTool) getMapping(filterIndex string, filterTypes string) eslib.MappingResponse {

	url := fmt.Sprintf("%s/_mapping/%s", filterIndex, filterTypes)
	request, err := (*tool.esClient).NewRequest("GET", url, "")
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

func (tool *MappingTool) GetIndexTypeAndTemplateLink() []TemplateLink {
	var links []TemplateLink

	// Parcours des index
	for indexName, indexValue := range tool.mappings {

		// Puis des types
		for typeName, mapping := range indexValue.MappingsByType {

			expectedMapping := tool.calculateExpectedMapping(indexName, typeName)
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

func (tool *MappingTool) calculateExpectedMapping(indexName string, typeName string) (expectedMapping eslib.MappingsJson) {
	var applicableTemplates []applicableTemplate = tool.searchForApplicableTemplate(indexName)
	sort.Sort(byOrderSort(applicableTemplates))

	for _, templateValue := range applicableTemplates {
		mergo.Map(&expectedMapping, templateValue.template.MappingsByType[typeName])
	}
	return expectedMapping
}

func (tool *MappingTool) searchForApplicableTemplate(indexName string) (result []applicableTemplate) {
	// Pour chacun des templates
	for templateName, templateValue := range tool.Templates {
		pattern := tool.getRegexPatternFromTemplateValue(templateValue.TemplateIndexPattern)
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

func (tool *MappingTool) getRegexPatternFromTemplateValue(templateIndexPattern string) (pattern string) {
	pattern = strings.Replace(templateIndexPattern, ".", "\\.", -1)
	pattern = strings.Replace(pattern, "*", ".*", -1)
	return pattern
}
