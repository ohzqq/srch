package srch

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

func FieldMappings(fields []string) []*mapping.FieldMapping {
	mapping := make([]*mapping.FieldMapping, len(fields))
	for i, field := range fields {
		mapping[i] = bleve.NewTextFieldMapping()
		mapping[i].Analyzer = StandardAnalyzer
		mapping[i].SkipFreqNorm = true
		mapping[i].DocValues = true
		mapping[i].Name = field
	}
	return mapping
}

func DocMapping(fields []*mapping.FieldMapping) *mapping.IndexMappingImpl {
	doc := bleve.NewIndexMapping()
	for _, field := range fields {
		doc.AddFieldMappingAt(field.Name, field)
	}
	return doc
}
