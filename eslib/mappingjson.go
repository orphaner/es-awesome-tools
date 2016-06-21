package eslib

type (
	MappingResponse map[string]MappingJson

	MappingJson struct {
		MappingsByType map[string]MappingsJson `json:"mappings"`
	}
)
