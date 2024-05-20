package srch

import (
	"testing"
)

type clientTest struct {
	*Client
}

func TestClientMem(t *testing.T) {
}

func TestClientDisk(t *testing.T) {
}

func TestClientNet(t *testing.T) {
}

func getTestClient(idx int) *Client {
	cfg := getTestCfg(idx)
	client, _ := NewClient(cfg)
	return client
}
