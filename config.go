package srch

type Config struct {
	Fields      []*Field `json:"fields"`
	Query       Query    `json:"filters"`
	Identifier  string   `json:"identifier"`
	interactive bool
	fuzzy       bool
}

func DefaultConfig() *Config {
	return &Config{
		Identifier: "id",
		Fields:     []*Field{NewTextField("title")},
	}
}
