package eslib

type (
	TemplateResponse map[string]TemplateJson

	TemplateJson struct {
		Order                int                     `json:"order"`
		TemplateIndexPattern string                  `json:"template"`
		MappingsByType       map[string]MappingsJson `json:"mappings"`
		Settings             struct {
			Index struct {
				RefreshInterval  string `json:"refresh_interval"`
				NumberOfReplicas int    `json:"number_of_replicas"`
				NumberOfShards   int    `json:"number_of_shards"`
			} `json:"index"`
		} `json:"settings"`
	}
)
