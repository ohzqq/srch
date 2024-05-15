package param

type Cfg struct {
	// Index Settings
	SrchAttr     []string `query:"searchableAttributes" json:"searchableAttributes,omitempty" mapstructure:"searchable_attributes" qs:"searchableAttributes"`
	FacetAttr    []string `query:"attributesForFaceting,omitempty" json:"attributesForFaceting,omitempty" mapstructure:"attributes_for_faceting" qs:"attributesForFaceting"`
	SortAttr     []string `query:"sortableAttributes,omitempty" json:"sortableAttributes,omitempty" mapstructure:"sortable_attributes" qs:"sortableAttributes"`
	DefaultField string   `query:"defaultField,omitempty" json:"defaultField,omitempty" mapstructure:"defaultField" qs:"defaultField"`
	UID          string   `query:"uid,omitempty" json:"uid,omitempty" mapstructure:"uid" qs:"uid"`
	Index        string   `query:"index,omitempty" json:"index,omitempty" mapstructure:"index" qs:"index"`
	Path         string   `json:"-" mapstructure:"path" qs:"path"`
}
