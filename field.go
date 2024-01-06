package srch

type Field struct {
	Attribute string           `json:"attribute"`
	Items     map[string][]int `json:"items,omitempty"`
	FullText  bool             `json:"fullText"`
}

func NewField(attr string) *Field {
	return &Field{
		Attribute: attr,
	}
}

func (f *Field) Add(value string, ids ...int) {
}

func (f *Field) addFullText(text string, ids ...int) {
	for _, token := range Tokenizer(text) {
		f.addTerm(token, ids...)
	}
}

func (f *Field) addTerm(term string, ids ...int) {
	if items, ok := f[term]; ok {
		items = append(items, ids...)
		return
	}
	f[term] = ids
}
