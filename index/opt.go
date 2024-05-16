package index

type Opt func(*Client) error

//func WithCfg(v any) Opt {
//  return func(idx *Client) error {
//    err := param.Decode(v, idx.Params)
//    if err != nil {
//      return err
//    }
//    return nil
//  }
//}

func WithRam(idx *Client) error {
	return idx.memDB()
}
