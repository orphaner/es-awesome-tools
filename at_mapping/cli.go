package at_mapping

import (
	"fmt"
	"github.com/orphaner/es-awesome-tools/eslib"
	"os"
	"sort"
	"text/tabwriter"
)

var w *tabwriter.Writer

func CliRun() {
	esClient := eslib.NewEsClient()
	esClient.SetFromFlag(Flags.Hostname)

	FillInData(esClient, Flags.Index, Flags.Types)
	printIndexTemplateSummary()
	printIndexMapping()
}

func printIndexTemplateSummary() {
	w = new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', tabwriter.Debug)
	for _, link := range GetIndexTypeAndTemplateLink() {
		templateName := link.TemplateName
		if link.TemplateName == "" {
			templateName = "NO TEMPLATE FOUND"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", link.IndexName, link.TypeName, templateName)
	}
	w.Flush()
}

func printIndexMapping() {
	for _, link := range GetIndexTypeAndTemplateLink() {
		fmt.Fprintf(os.Stdout, "----- %s/%s -----\n", link.IndexName, link.TypeName)

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", "Property",
			"Mapping Type", "Template Type",
			"Mapping Store", "Template Store",
			"Mapping Index", "Template Index",
			"Mapping Format", "Template Format")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", "________",
			"____________", "_____________",
			"_____________", "______________",
			"_____________", "______________",
			"______________", "_______________")

		writeProperties("", link.EffectiveMapping.Properties, link.ExpectedMapping.Properties)
		w.Flush()
		fmt.Fprintf(os.Stdout, "\n")
	}
}

func writeProperties(levelName string, effectiveProperties map[string]eslib.PropertyJson, expectedProperties map[string]eslib.PropertyJson) {
	sortedKeys := getSortedKeys(effectiveProperties)

	for _, propName := range sortedKeys {
		effectiveProp := effectiveProperties[propName]
		expectedProp := expectedProperties[propName]

		displayPropName := propName
		if levelName != "" {
			displayPropName = fmt.Sprintf("%s.%s", levelName, propName)
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", displayPropName,
			effectiveProp.Type, expectedProp.Type,
			effectiveProp.Store, expectedProp.Store,
			effectiveProp.Index, expectedProp.Index,
			effectiveProp.Format, expectedProp.Format)

		if effectiveProp.Type == "nested" {
			writeProperties(propName, effectiveProp.Properties, expectedProp.Properties)
		}
		if effectiveProp.Fields != nil {
			writeProperties(propName, effectiveProp.Fields, expectedProp.Fields)
		}
	}
}

func getSortedKeys(indexMapping map[string]eslib.PropertyJson) []string {
	mapSize := len(indexMapping)
	keys := make([]string, mapSize)
	i := 0
	for key, _ := range indexMapping {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}
