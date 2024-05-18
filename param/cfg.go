package param

type Cfg struct {
	Client *Client
	Search *Search
	Idx    *Idx
}

func NewCfg() *Cfg {
	return &Cfg{
		Idx:    NewIdx(),
		Search: NewSearch(),
		Client: NewClient(),
	}
}

func (cfg *Cfg) Decode(v any) error {
	err := Decode(v, cfg.Idx)
	if err != nil {
		return err
	}
	err = Decode(v, cfg.Search)
	if err != nil {
		return err
	}
	err = Decode(v, cfg.Client)
	if err != nil {
		return err
	}
	return nil
}
