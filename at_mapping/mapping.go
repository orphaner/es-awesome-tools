package at_mapping

import (
	"github.com/orphaner/es-awesome-tools/eslib"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type (
	TemplateLink struct {
		IndexName    string
		TypeName     string
		TemplateName string
		Mapping      *eslib.MappingsJson
		Template     *eslib.TemplateJson
	}
	byIndexAndTypeSort []*TemplateLink
)

var (
	esClient  *eslib.EsClient
	Templates *eslib.TemplateResponse
	Mappings  *eslib.MappingResponse
	Links     []*TemplateLink
)

func (by byIndexAndTypeSort) Len() int {
	return len(by)
}
func (by byIndexAndTypeSort) Swap(i, j int) {
	by[i], by[j] = by[j], by[i]
}
func (by byIndexAndTypeSort) Less(i, j int) bool {
	if strings.Compare(by[i].IndexName, by[j].IndexName) == -1 {
		return true
	}
	if strings.Compare(by[i].TypeName, by[i].TypeName) == -1 {
		return true
	}
	return false
}

func FillInData(esClientParam *eslib.EsClient, filterIndex string, filterTypes string) {
	esClient = esClientParam

	Templates = getTemplate()
	Mappings = getMapping(filterIndex, filterTypes)
	Links = linkIndexTypeAndTemplate()
	sort.Sort(byIndexAndTypeSort(Links))
}

func getTemplate() *eslib.TemplateResponse {

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

	return &template
}

func getMapping(filterIndex string, filterTypes string) *eslib.MappingResponse {

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

	return &mapping
}

func linkIndexTypeAndTemplate() []*TemplateLink {
	var links []*TemplateLink

	// Parcours des index
	for indexName, indexValue := range *Mappings {

		// Puis des types
		for typeName, mapping := range indexValue.MappingsByType {

			templateName, template := searchForTemplate(indexName)
			links = append(links, &TemplateLink{
				IndexName:    indexName,
				TypeName:     typeName,
				TemplateName: templateName,
				Mapping:      &mapping,
				Template:     template,
			})
		}
	}
	return links
}

func searchForTemplate(indexName string) (string, *eslib.TemplateJson) {
	for templateName, templateValue := range *Templates {
		pattern := getRegexPatternFromTemplateValue(templateValue.TemplateIndexPattern)
		regex := regexp.MustCompile(pattern)
		if regex.MatchString(indexName) {
			return templateName, &templateValue
		}
	}
	return "", nil
}

func getRegexPatternFromTemplateValue(templateIndexPattern string) (pattern string) {
	pattern = strings.Replace(templateIndexPattern, ".", "\\.", -1)
	pattern = strings.Replace(pattern, "*", ".*", -1)
	return pattern
}
