package eslib

type (
	MappingResponse map[string] struct {
		MappingsByType map[string]MappingsJson `json:"mappings"`
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
