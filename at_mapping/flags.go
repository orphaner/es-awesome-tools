package at_mapping

import "flag"

type (
	VarFlag struct {
		Hostname string
		Index    string
		Types    string
	}
)

var (
	Flags VarFlag
)

func init() {
	flag.StringVar(&Flags.Hostname, "host", "http://127.0.0.1:9200", "Elasticsearch hostname/ip and port ; pattern: http[s]://[user[:pass]@]ip:port ")
	flag.StringVar(&Flags.Index, "index", "*", "Elasticsearch indexes to analyze, separated by comma ; wildcard allowed")
	flag.StringVar(&Flags.Types, "types", "*", "Elasticsearch types to analyze, separated by comma ; wildcard allowed")
	flag.Parse()
}
