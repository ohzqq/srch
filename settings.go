package srch

import "strings"

type Settings struct {
	SearchableAttributes  []string
	AttributesForFaceting []string
}

func NewSettings(q any) *Settings {
	settings := &Settings{
		SearchableAttributes: []string{"title"},
	}

	v, err := ParseValues(q)
	if err != nil {
		return settings
	}

	for k, _ := range v {
		switch k {
		case "searchableAttributes":
			settings.SearchableAttributes = strings.Split(v.Get(k), ",")
		case "attributesForFaceting":
			settings.AttributesForFaceting = strings.Split(v.Get(k), ",")
		}
	}
	return settings
}
