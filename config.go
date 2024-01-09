package srch

type Config struct {
	Fields      []*Field `json:"fields"`
	Query       Query    `json:"filters"`
	Identifier  string   `json:"identifier"`
	interactive bool
	fuzzy       bool
}
