package param

import (
	"net/url"

	"github.com/spf13/viper"
)

type Search struct {
	Hits                 int      `query:"hits,omitempty" json:"hits,omitempty"`
	AttributesToRetrieve []string `query:"attributesToRetrieve,omitempty" json:"attributesToRetrieve,omitempty"`
	Page                 int      `query:"page,omitempty" json:"page,omitempty"`
	HitsPerPage          int      `query:"hitsPerPage,omitempty" json:"hitsPerPage,omitempty"`
	Query                string   `query:"query,omitempty" json:"query,omitempty"`
	SortBy               string   `query:"sortBy,omitempty" json:"sortBy,omitempty"`
	Order                string   `query:"order,omitempty" json:"order,omitempty"`
	*FacetSettings
}

func NewSearch() *Search {
	return &Search{
		HitsPerPage:   viper.GetInt("hitsPerPage"),
		FacetSettings: NewFacetSettings(),
	}
}

func (s *Search) Parse(q string) error {
	vals, err := url.ParseQuery(q)
	if err != nil {
		return err
	}
	return s.Set(vals)
}

func (s *Search) Set(v url.Values) error {
	for _, key := range paramsSearch {
		switch key {
		case Hits:
			s.Hits = GetQueryInt(key, v)
		case AttributesToRetrieve:
			s.AttributesToRetrieve = GetQueryStringSlice(key, v)
		case Page:
			s.Page = GetQueryInt(key, v)
		case HitsPerPage:
			s.HitsPerPage = GetQueryInt(key, v)
		case Query:
			s.Query = v.Get(key)
		case SortBy:
			s.SortBy = v.Get(key)
		case Order:
			s.Order = v.Get(key)
		}
		v.Del(key)
	}
	s.FacetSettings.Set(v)
	return nil
}

func (s *Search) Has(key string) bool {
	switch key {
	case Hits:
		return s.Hits != 0
	case AttributesToRetrieve:
		return len(s.AttributesToRetrieve) != 0
	case Page:
		return s.Page != 0
	case HitsPerPage:
		return s.HitsPerPage != 0
	case Query:
		return s.Query != ""
	case SortBy:
		return s.SortBy != ""
	case Order:
		return s.Order != ""
	}
	return s.FacetSettings.Has(key)
}
