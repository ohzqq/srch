package index

type Opt func(*Client) error

func WithRam(idx *Client) error {
	return idx.memDB()
}
