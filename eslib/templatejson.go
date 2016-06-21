package eslib

type (
	TemplateResponse map[string]TemplateJson

	TemplateJson struct {
		Order          int                     `json:"order"`
		Template       string                  `json:"template"`
		MappingsByType map[string]MappingsJson `json:"mappings"`
		Settings       struct {
			Index struct {
				RefreshInterval  string `json:"refresh_interval"`
				NumberOfReplicas int    `json:"number_of_replicas"`
				NumberOfShards   int    `json:"number_of_shards"`
			} `json:"index"`
		} `json:"settings"`
	}
	MappingsJson struct {
		Properties map[string]PropertyJson `json:"properties"`
	}
	PropertyJson struct {
		Type       string                  `json:"type"`
		Format     string                  `json:"format"`
		Index      string                  `json:"index"`
		Store      string                  `json:"store"`
		Fields     map[string]PropertyJson `json:"fields"`
		Properties map[string]PropertyJson `json:"properties"`
	}
)
