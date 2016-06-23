package at_mapping

import (
	"github.com/orphaner/es-awesome-tools/eslib"
	"strings"
)

type (
	byIndexAndTypeSort []TemplateLink
	byOrderSort        []eslib.TemplateJson
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


func (by byOrderSort) Len() int {
	return len(by)
}
func (by byOrderSort) Swap(i, j int) {
	by[i], by[j] = by[j], by[i]
}
func (by byOrderSort) Less(i, j int) bool {
	if by[i].Order < by[j].Order {
		return true
	}
	return false
}
