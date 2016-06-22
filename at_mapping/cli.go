package at_mapping

import (
	"github.com/orphaner/es-awesome-tools/eslib"
	"fmt"
	"os"
	"text/tabwriter"
)

func CliRun() {
	esClient := eslib.NewEsClient()
	esClient.SetFromFlag(Flags.Hostname)

	FillInData(esClient, Flags.Index, Flags.Types)
	writeIndexTemplateSummary()
}

func writeIndexTemplateSummary() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	for _, link := range Links {
		templateName := link.TemplateName
		if link.TemplateName == "" {
			templateName = "NO TEMPLATE FOUND"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", link.IndexName, link.TypeName, templateName)
	}
	w.Flush()
}
