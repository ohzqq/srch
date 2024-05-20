package srch

import "testing"

func TestClientMem(t *testing.T) {
}

func TestClientDisk(t *testing.T) {
}

func TestClientNet(t *testing.T) {
}

//func TestClientInitStr(t *testing.T) {
//  for query, test := range ParamTests() {
//    client, err := NewClient(query.String())
//    if err != nil {
//      t.Fatal(err)
//    }

//    //if cfg.DB != "" {
//    //  if cfg.Path != hareTestPath {
//    //    t.Errorf("got %v, wanted %v\n", cfg.Path, hareTestPath)
//    //  }
//    //  if cfg.Scheme != "file" {
//    //    t.Errorf("got %v, wanted %v\n", cfg.Scheme, "file")
//    //  }
//    //}

//    if !client.TableExists(settingsTbl) {
//      t.Error(test.Err("", errors.New("_settings table doesn't exist")))
//    }
//    _, err = client.GetIdxCfg(client.Client.Index)
//    if err != nil {
//      t.Error(test.Err(test.Msg("", client.Client.Index, test.Cfg.IndexName()), err))
//    }
//  }
//}
